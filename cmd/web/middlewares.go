package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self' https://unpkg.com/htmx.org@1.9.12 https://unpkg.com/htmx.org@1.9.12/dist/ext/json-enc.js https://unpkg.com/htmx.org@1.9.12/dist/ext/response-targets.js; style-src 'self' fonts.googleapis.com; font-src 'self' fonts.gstatic.com data:")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		
		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip = r.RemoteAddr
			proto = r.Proto
			method = r.Method
			uri = r.URL.RequestURI()
		)
		logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var (
					method = r.Method
					uri = r.URL.RequestURI()
					trace = string(debug.Stack())
				)
				w.Header().Set("Connection", "close")
				slog.Error("Internal server error", "method", method, "path", uri, "trace", trace)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}			
		}()

		next.ServeHTTP(w, r)
	})
}