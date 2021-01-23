package database

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var fireStore *firestore.Client

func Connect() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fs, err := app.Firestore(context.Background())
	if err != nil {
		panic(err)
	}
	fireStore = fs
}

func Sessions() *firestore.CollectionRef {
	return fireStore.Collection("sessions")
}
