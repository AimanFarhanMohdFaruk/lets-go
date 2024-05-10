package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ApiError struct {
	StatusCode int "json:\"statusCode\""
	Error string "json:\"msg\""
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

func ServerError(w http.ResponseWriter, r *http.Request, status int, err error) {
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "path", uri, "trace", trace)
	WriteJSON(w, status, ApiError{StatusCode: status, Error: err.Error()})
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}