package main

import (
	"assignment2/APIs"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handlerRequests() {
	http.HandleFunc("/repocheck/v1/commits", APIs.CommitsHandler)
	http.HandleFunc("/repocheck/v1/languages", APIs.LanguagesHandler)
	http.HandleFunc("/repocheck/v1/webhooks/", APIs.WebhookWithIdHandler)
	http.HandleFunc("/repocheck/v1/webhooks", APIs.WebhookHandler)
	http.HandleFunc("/repocheck/v1/status", APIs.StatusHandler)
}

func main() {

	fmt.Println("heey")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	handlerRequests()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
