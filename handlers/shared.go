package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	StatusCode int "json:\"statusCode\""
	Msg any "json:\"msg\""
}

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w,r); err != nil {
			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json") // important, must set the header to json before writing the status
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func RenderError(w http.ResponseWriter, r *http.Request, status int, err error){
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "path", uri, "trace", trace)
	WriteJSON(w, status, ApiError{
		StatusCode: status,
		Msg: err.Error(),
	})
}

func InvalidRequestData(w http.ResponseWriter, r *http.Request, err error) {
	fieldErrors := GetFieldErrors(err)

	var (
		method = r.Method
		uri = r.URL.RequestURI()
	)

	slog.Error(err.Error(), "method", method, "path", uri)
	WriteJSON(w, http.StatusUnprocessableEntity, ApiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg: fieldErrors,
	})
}

func GetFieldErrors(err error) map[string]string {
	fieldErrors := make(map[string] string)
	for _ , err := range err.(validator.ValidationErrors) {
		fieldErrors[err.Field()] = err.Tag()
	}

	return fieldErrors
}

// Function below is used for errors that we do not want to show to the user.
// the full error will get logged
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "path", uri, "trace", trace)
	WriteJSON(w, http.StatusInternalServerError,
		ApiError{
			StatusCode: http.StatusInternalServerError,
		 	Msg: http.StatusText(http.StatusInternalServerError)})
}
