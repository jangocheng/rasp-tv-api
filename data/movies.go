package data

import (
	"database/sql"
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
