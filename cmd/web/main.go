package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // $ symbol restructs this route to exact match on / only
	mux.HandleFunc("GET /snippet/view/{id}", showSnippet)
	mux.HandleFunc("GET /snippet/create", newSnippetForm)
	mux.HandleFunc("POST /snippet/create", createSnippet)

	log.Print("Starting server on :3000")

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}