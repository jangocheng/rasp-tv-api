package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
)

type Movie struct {
	Id        int64
	Title     sql.NullString
	Filepath  string
	Length    sql.NullFloat64
	IsIndexed bool
}

func (m *Movie) MarshalJSON() ([]byte, error) {
	var title *string
	if m.Title.Valid {
		title = &m.Title.String
	} else {
		title = nil
	}

	var length *float64
	if m.Length.Valid {
		length = &m.Length.Float64
	} else {
		length = nil
	}

	mov := struct {
		Id        int64    `json:"id"`
		Title     *string  `json:"title"`
		Filepath  string   `json:"filepath"`
		Length    *float64 `json:"length"`
		IsIndexed bool     `json:"isIndexed"`
	}{
		m.Id,
		title,
		m.Filepath,
		length,
		m.IsIndexed,
	}

	return json.Marshal(&mov)
}

func (m *Movie) UnmarshalJSON(data []byte) error {
	var mov struct {
		Id        int64
		Title     *string
		Filepath  string
		Length    *float64
		IsIndexed bool
	}

	if err := json.Unmarshal(data, &mov); err != nil {
		return err
	}

	m.Id = mov.Id
	m.Filepath = mov.Filepath
	m.IsIndexed = mov.IsIndexed

	if mov.Title == nil {
		m.Title = sql.NullString{Valid: false}
	} else {
		m.Title = sql.NullString{Valid: true, String: *mov.Title}
	}

	if mov.Length == nil {
		m.Length = sql.NullFloat64{Valid: false}
	} else {
		m.Length = sql.NullFloat64{Valid: true, Float64: *mov.Length}
	}

	return nil
}

func (m *Movie) Update(db *sql.DB) error {
	if !m.Title.Valid {
		return fmt.Errorf("Cannot update movie with invalid title")
	}

	_, err := db.Exec("UPDATE movies SET title = ?, isIndexed = 1 WHERE id = ?", m.Title, m.Id)
	return err
}

func (m *Movie) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM movies WHERE Id = ?", m.Id)
	return err
}

func GetMovies(filter string, db *sql.DB) ([]Movie, error) {
	movies := make([]Movie, 0, 70)
	rows, err := db.Query("SELECT id, title, filepath, length, isIndexed FROM movies " + filter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return movies, nil
	}

	for rows.Next() {
		m := Movie{}
		if err := rows.Scan(&m.Id, &m.Title, &m.Filepath, &m.Length, &m.IsIndexed); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}
