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

	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/signup", SignUp)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
