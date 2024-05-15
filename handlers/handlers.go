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

func ShowHomePage(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Index()
		if err != nil {
			ServerError(w, r, err)
		}
		
		component := ui.Home(snippets)
		err = component.Render(r.Context(), w)
		if err != nil {
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
			http.NotFound(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
			RenderError(w, r, http.StatusNotFound, err)
		} else {
			ServerError(w, r, err)
		}
			return
		}

		page := ui.View(snippet)
		err = page.Render(r.Context(), w)
		if err != nil {
			ServerError(w, r, err)
		}
	})
}

func NewSnippetForm(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form := ui.SnippetCreateForm{
			Title: "",
			Content: "",
			Expires: 365,
			FieldErrors: make(map[string]string),
		}

		page := ui.Create(form)

		err := page.Render(r.Context(), w)
		if err != nil {
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
		
		fieldErrors := make(map[string]string)
		if err := app.Validator.Struct(createSnippetRequest); err != nil {
			fieldErrors = GetFieldErrors(err)
		}

		if len(fieldErrors) > 0 {
			form := ui.SnippetCreateForm{
				Title: createSnippetRequest.Title,
				Content: createSnippetRequest.Content,
				Expires: createSnippetRequest.Expires,
				FieldErrors: fieldErrors,
			}
			page := ui.Create(form)
			err := page.Render(r.Context(), w)

			if err != nil {
				ServerError(w, r, err)
			}
			return
		}

		id, err := app.Snippets.Create(createSnippetRequest.Title, createSnippetRequest.Content, createSnippetRequest.Expires)
		if err != nil {
			ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	})
}

func GetLatestSnippets(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Index()
		if err != nil {
			ServerError(w, r, err)
		}

		WriteJSON(w, http.StatusOK, snippets)
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