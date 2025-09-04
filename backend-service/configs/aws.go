package configs

import (
	"context"
	"errors"
	"jirbthagoras/raksana-backend/helpers"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type AWSClient struct {
	S3Client          *s3.Client
	RekognitionClient *rekognition.Client
	PsClient          *s3.PresignClient
}

func InitAWSClient(cnf *viper.Viper) *AWSClient {
	region := cnf.GetString("AWS_REGION")
	accessKey := cnf.GetString("AWS_ACCESS_KEY_ID")
	secretKey := cnf.GetString("AWS_SECRET_ACCESS_KEY")
	sessionToken := cnf.GetString("AWS_SESSION_TOKEN")

	customCredentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey,
		secretKey,
		sessionToken,
	))

	awsCnf, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(region),
		awsConfig.WithCredentialsProvider(customCredentials),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(awsCnf)
	psClient := s3.NewPresignClient(s3Client)
	rekognitionClient := rekognition.NewFromConfig(awsCnf)

	slog.Debug("Established connection to AWS")
	return &AWSClient{
		S3Client:          s3Client,
		RekognitionClient: rekognitionClient,
		PsClient:          psClient,
	}
}

func (a *AWSClient) CreatePresignUrlPutObject(key string, contentType string) (*v4.PresignedHTTPRequest, error) {
	cnf := helpers.NewConfig()
	bucket := cnf.GetString("AWS_BUCKET")

	presignReq, err := a.PsClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(10*time.Minute))
	if err != nil {
		slog.Error("Failed to create a presigned url")
		return nil, err
	}

	return presignReq, nil
}

func (a *AWSClient) DeleteObject(bucketName string, key string) error {
	ctx := context.Background()
	_, err := a.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("Failed to delete image from bucket")
		return err
	}
	return nil
}

func (a *AWSClient) CheckObjectExistence(bucketName string, key string) error {
	ctx := context.Background()
	_, err := a.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError

		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return fiber.NewError(fiber.StatusBadRequest, "No pofile picture uploaded")
			}
			if respErr, ok := err.(interface{ HTTPStatusCode() int }); ok && respErr.HTTPStatusCode() == http.StatusNotFound {
				slog.Error("Failed to check object head")
				return err
			}
		}

		slog.Error("Failed to check object head")
		return err
	}

	return nil
}
