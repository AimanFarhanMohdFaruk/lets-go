package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aiman-farhan/snippetbox/config"
	"github.com/aiman-farhan/snippetbox/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("LISTEN_ADDR")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))
	logger.Info("Starting server", "addr" , addr)

	db, err := config.NewDB()
	if err != nil {
		logger.Error("Error connecting to database", err, err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &config.Application{
		DB: db,
		Logger: logger,
		Snippets: &models.SnippetModel{
			DB: db,
		},
	}
	
	err = http.ListenAndServe(addr, routes(app))
	if err != nil {
		logger.Error("Error starting server", err, err.Error())
	}

	os.Exit(1)
}