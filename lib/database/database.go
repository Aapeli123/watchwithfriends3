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

	// Clear all old sessions and rooms
	sessions, err := fs.Collection("sessions").Documents(context.Background()).GetAll()
	if err != nil {
		panic(err)
	}
	for _, s := range sessions {
		s.Ref.Delete(context.Background())
	}

	rooms, err := fs.Collection("rooms").Documents(context.Background()).GetAll()
	if err != nil {
		panic(err)
	}
	for _, r := range rooms {
		r.Ref.Delete(context.Background())
	}

	fireStore = fs
}

func Sessions() *firestore.CollectionRef {
	return fireStore.Collection("sessions")
}

func Rooms() *firestore.CollectionRef {
	return fireStore.Collection("rooms")
}
