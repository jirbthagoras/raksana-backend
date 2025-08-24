package app

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

type AWSClient struct {
	S3Client          *s3.Client
	RekognitionClient *rekognition.Client
	Uploader          *manager.Uploader
}

func InitS3Client(cnf *viper.Viper) *AWSClient {
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
	uploader := manager.NewUploader(s3Client)
	rekognitionClient := rekognition.NewFromConfig(awsCnf)

	return &AWSClient{
		S3Client:          s3Client,
		RekognitionClient: rekognitionClient,
		Uploader:          uploader,
	}
}
