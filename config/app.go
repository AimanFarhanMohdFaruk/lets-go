package config

import (
	"database/sql"
	"log/slog"

	"github.com/aiman-farhan/snippetbox/internal/models"
	"github.com/go-playground/validator/v10"
)

type Application struct {
	DB *sql.DB
	Logger *slog.Logger
	Snippets *models.SnippetModel
	Validator *validator.Validate
}