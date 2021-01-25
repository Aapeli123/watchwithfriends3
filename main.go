package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Aapeli123/watchwithfriends3/lib/database"
	"github.com/Aapeli123/watchwithfriends3/server"
)

func main() {
	rand.Seed(time.Now().Unix()) // Seed random generator for the room codes
	database.Connect()
	fmt.Println("Watchwithfriends Server v3 starting...")
	server.Start(server.Options{
		Address: ":8080",
		HTTPS:   false,
	})
}
