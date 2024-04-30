package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	ui "github.com/aiman-farhan/snippetbox/ui/html/pages"
)

func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	component := ui.Home()

	err := component.Render(r.Context(), w)
	if err != nil {
		log.Print(err.Error())
		ServerError(w, r, err)
		return
	}
}	

func ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with id: %d", id)
}

func NewSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet ...."))
}
