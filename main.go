package main

import (
	"diocforum/auth"
	"diocforum/post"
	"diocforum/thread"
	"log"
	"net/http"
)

var apiBase = "/api"

func main() {
	thread.SetupRoute(apiBase)
	post.SetupRoute(apiBase)
	auth.SetUpRoute(apiBase)

	log.Println("Server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
