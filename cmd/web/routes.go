package main

import (
	"net/http"

	"github.com/aiman-farhan/snippetbox/config"
	"github.com/aiman-farhan/snippetbox/handlers"
)

func routes(app *config.Application) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", handlers.ShowHomePage(app)) // $ symbol restructs this route to exact match on / only
	mux.Handle("GET /snippets/view/{id}", handlers.ShowSnippet(app))
	mux.Handle("GET /snippets/create", handlers.NewSnippetForm(app))
	mux.Handle("POST /snippets/create", handlers.CreateSnippet(app))

	return recoverPanic(logRequest(commonHeaders(mux), app.Logger))
}