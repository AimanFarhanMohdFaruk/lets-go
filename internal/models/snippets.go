package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

type Snippet struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expire int) (int, error){
	stmt := `INSERT INTO snippets 
	(title, content, created, expires) 
	VALUES ($1, $2, $3, $4) RETURNING id`

	var lastInsertedId int64 = 0
	err := m.DB.QueryRow(stmt,title, content, time.Now().UTC(), time.Now().Add(time.Hour * 24 * 7).UTC()).Scan(&lastInsertedId)
	
	if err != nil {
		slog.Error("Error creating snippet", err, err.Error() )
	}

	return int(lastInsertedId), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP ORDER BY ID DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
