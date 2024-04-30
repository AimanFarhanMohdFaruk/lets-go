package main

import (
	"net/http"

	"github.com/aiman-farhan/snippetbox/handlers"
)
	

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", handlers.ShowHomePage) // $ symbol restructs this route to exact match on / only
	mux.HandleFunc("GET /snippet/view/{id}", handlers.ShowSnippet)
	mux.HandleFunc("GET /snippet/create", handlers.NewSnippetForm)
	mux.HandleFunc("POST /snippet/create", handlers.CreateSnippet)

	return mux
}