package configs

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"log/slog"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func InitFCMClient() *messaging.Client {
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	cnf := helpers.NewConfig()

	firebaseConfig := &firebase.Config{
		ProjectID: cnf.GetString("FIREBASE_PROJECT_ID"),
	}

	app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		slog.Error("Failed to connect to firebase app", "err", err.Error())
		os.Exit(1)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		slog.Error("Failed to get FCM Client", "err", err.Error())
		os.Exit(1)
	}

	slog.Debug("Established connection to FCM")
	return client
}
