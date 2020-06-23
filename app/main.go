package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
