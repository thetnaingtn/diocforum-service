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

	log.Println("Server listening at :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
