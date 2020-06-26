package main

import (
	"github.com/gorilla/handlers"
	"gomail/database"
	"gomail/email"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	err = email.ConnectToRabbitMQ()
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	router := newRouter()
	log.Print("\nServer started on port 8080")
	// start listening to port 8080
	err = http.ListenAndServe(
		":8080",
		handlers.CORS(
			// Allowed origins are specified in docker-compose.yaml
			handlers.AllowedOrigins(strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
		)(router),
	)

}
