package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/ path")
	})

	log.Fatal(http.ListenAndServe(":5000", nil))
}
