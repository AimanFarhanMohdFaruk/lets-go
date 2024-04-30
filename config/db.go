package config

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type Env struct {
	Db *sql.DB
}

func NewDB() (*sql.DB, error) {
	connStr := "user=weblocalhost dbname=snippetbox password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	slog.Info("Database snippetbox connected!")

	return db, nil
}