package main

import (
	"gomail/database"
	"log"
	"net/http"
)

func main() {

	err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	router := newRouter()

	log.Print("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
