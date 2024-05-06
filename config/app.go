package config

import (
	"database/sql"
	"log/slog"

	"github.com/aiman-farhan/snippetbox/internal/models"
)

type Application struct {
	DB *sql.DB
	Logger *slog.Logger
	Snippets *models.SnippetModel
}