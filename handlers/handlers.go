package handlers

import (
	"encoding/json"
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
		snippets, err := app.Snippets.Index()
		if err != nil {
			ServerError(w, r, err)
		}

		WriteJSON(w, http.StatusOK, snippets)
	})
}

func ShowHomePage(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Index()
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
	})
}

func NewSnippetForm(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := ui.Create()

		err := page.Render(r.Context(), w)
		if err != nil {
			app.Logger.Error(err.Error())
			ServerError(w, r, err)
		}
		app.Logger.Info("Rendering create form")
	})
}
 
func CreateSnippet(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createSnippetRequest := models.CreateSnippetRequest{}
		
		if err := json.NewDecoder(r.Body).Decode(&createSnippetRequest); err != nil {
			ServerError(w, r, err)
			return
		}
		defer r.Body.Close()

		if err := app.Validator.Struct(createSnippetRequest); err != nil {
			InvalidRequestData(w, r, err)
		}

		id, err := app.Snippets.Create(createSnippetRequest.Title, createSnippetRequest.Content, createSnippetRequest.Expires)
		if err != nil {
			ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	})
}

// Pattern for including a panic recovery for any background jobs

// func (app *application) myHandler(w http.ResponseWriter, r *http.Request) {
// 	// Spin up a new goroutine to do some background processing.
// 	go func() {
// 	defer func() {
// 	if err := recover(); err != nil {
// 	app.logger.Error(fmt.Sprint(err))
// 	}
// 	}()
// 	doSomeBackgroundProcessing()
// 	}()
// 	w.Write([]byte("OK"))
// }