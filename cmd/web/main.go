package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", ":3000", "HTTP Network Address")
	flag.Parse()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", showHomePage) // $ symbol restructs this route to exact match on / only
	
	mux.HandleFunc("GET /snippet/view/{id}", showSnippet)
	mux.HandleFunc("GET /snippet/create", newSnippetForm)
	mux.HandleFunc("POST /snippet/create", createSnippet)

	log.Printf("Starting server on port %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}