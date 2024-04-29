package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aiman-farhan/snippetbox/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("PORT")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", handlers.ShowHomePage) // $ symbol restructs this route to exact match on / only
	
	mux.HandleFunc("GET /snippet/view/{id}", handlers.ShowSnippet)
	mux.HandleFunc("GET /snippet/create", handlers.NewSnippetForm)
	mux.HandleFunc("POST /snippet/create", handlers.CreateSnippet)

	logger.Info("Starting on server", "addr" , addr)

	start := http.ListenAndServe(addr, mux)
	logger.Error(start.Error())
	os.Exit(1)
}