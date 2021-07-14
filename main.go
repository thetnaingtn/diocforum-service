package main

import (
	"diocforum/auth"
	"diocforum/post"
	"log"
	"net/http"
)

var apiBase = "/api"

func main() {
	post.SetupRoute(apiBase)
	auth.SetUpRoute(apiBase)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
