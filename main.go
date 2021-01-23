package main

import (
	"fmt"

	"github.com/Aapeli123/watchwithfriends3/server"
)

func main() {
	fmt.Println("Watchwithfriends Server v3 starting...")
	server.Start(server.Options{
		Address: ":8080",
		HTTPS:   false,
	})
}
