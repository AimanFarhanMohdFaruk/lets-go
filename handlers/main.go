package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aiman-farhan/snippetbox/config"
	"github.com/aiman-farhan/snippetbox/internal/models"
	ui "github.com/aiman-farhan/snippetbox/ui/html/pages"
)

func GetLatestSnippets(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			ServerError(w, r, err)
		}

		WriteJSON(w, http.StatusOK, snippets)
	})
}

func ShowHomePage(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")
	
		snippets, err := app.Snippets.Latest()
		if err != nil {
			ServerError(w, r, err)
		}
		
		component := ui.Home(snippets)
		err = component.Render(r.Context(), w)
		if err != nil {
			app.Logger.Error(err.Error())
			ServerError(w, r, err)
			return
		}
		app.Logger.Info("Rendering home page")
	})
}

func ShowSnippet(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			app.Logger.Error("Invalid identifier for snippet")
			http.NotFound(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			ServerError(w, r, err)
		}
			return
		}

		component := ui.View(snippet)
		err = component.Render(r.Context(), w)
		if err != nil {
			app.Logger.Error(err.Error())
			ServerError(w, r, err)
		}
		app.Logger.Info("Rendering view snippet")
	})
}

func NewSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func CreateSnippet(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	})
}
