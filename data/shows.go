package data

import (
	"database/sql"
	"fmt"
	"io"
)

type Show struct {
	Id       int64
	Title    string
	Episodes []Episode
}

func (s *Show) Add(db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO shows (title) VALUES (?)", s.Title)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, err
}

type Episode struct {
	Id        int64
	ShowId    sql.NullInt64
	Title     sql.NullString
	Number    sql.NullInt64
	Season    sql.NullInt64
	Filepath  string
	Length    sql.NullFloat64
	IsIndexed bool
}

func (e *Episode) Update(db *sql.DB) error {
	if !e.ShowId.Valid {
		return fmt.Errorf("Cannot update episode with invalid showId")
	}

	if !e.Title.Valid {
		return fmt.Errorf("Cannot update episode with invalid title")
	}

	if !e.Number.Valid {
		return fmt.Errorf("Cannot update episode with invalid episode number")
	}

	if !e.Season.Valid {
		return fmt.Errorf("Cannot update episode with invalid season")
	}

	query := "UPDATE episodes SET showId = ?, title = ?, episodeNumber = ?, season = ?, isIndexed = 1 WHERE id = ?"
	_, err := db.Exec(query, e.ShowId, e.Title, e.Number, e.Season, e.Id)
	return err
}

func (e *Episode) DeleteEpisode(db *sql.DB) error {
	var err error
	_, err = db.Exec("DELETE FROM episodes WHERE Id = ?;", e.Id)
	return err
}

func GetShows(filter string, db *sql.DB) ([]Show, error) {
	shows := make([]Show, 0, 10)
	rows, err := db.Query("SELECT id, title FROM shows " + filter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return shows, nil
	}

	for rows.Next() {
		show := Show{}
		rows.Scan(&show.Id, &show.Title)
		shows = append(shows, show)
	}

	return shows, nil
}

func GetEpisodes(filter string, db *sql.DB) ([]Episode, error) {
	episodes := make([]Episode, 0, 20)
	rows, err := db.Query("SELECT id, title, episodeNumber, season, filepath, length, isIndexed, showId FROM episodes " + filter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return episodes, nil
	}

	for rows.Next() {
		e := Episode{}
		if err := rows.Scan(&e.Id, &e.Title, &e.Number, &e.Season, &e.Filepath, &e.Length, &e.IsIndexed, &e.ShowId); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}

	return episodes, nil
}
