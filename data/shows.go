package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
)

type Show struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
	Episodes []Episode `json:"episodes,omitempty"`
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
	s.Id = id

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

func (e *Episode) MarshalJSON() ([]byte, error) {
	var showId *int64
	if e.ShowId.Valid {
		showId = &e.ShowId.Int64
	} else {
		showId = nil
	}

	var title *string
	if e.Title.Valid {
		title = &e.Title.String
	} else {
		title = nil
	}

	var number *int64
	if e.Number.Valid {
		number = &e.Number.Int64
	} else {
		number = nil
	}

	var season *int64
	if e.Season.Valid {
		season = &e.Season.Int64
	} else {
		season = nil
	}

	var length *float64
	if e.Length.Valid {
		length = &e.Length.Float64
	} else {
		length = nil
	}

	episode := struct {
		Id        int64    `json:"id"`
		ShowId    *int64   `json:"showId"`
		Title     *string  `json:"title"`
		Number    *int64   `json:"number"`
		Season    *int64   `json:"season"`
		Filepath  string   `json:"filepath"`
		Length    *float64 `json:"length"`
		IsIndexed bool     `json:"isIndexed"`
	}{
		e.Id,
		showId,
		title,
		number,
		season,
		e.Filepath,
		length,
		e.IsIndexed,
	}

	return json.Marshal(&episode)
}

func (e *Episode) UnmarshalJSON(data []byte) error {
	var episode struct {
		Id        int64
		ShowId    *int64
		Title     *string
		Number    *int64
		Season    *int64
		Filepath  string
		Length    *float64
		IsIndexed bool
	}

	if err := json.Unmarshal(data, &episode); err != nil {
		return err
	}

	e.Id = episode.Id
	e.Filepath = episode.Filepath
	e.IsIndexed = episode.IsIndexed

	if episode.ShowId == nil {
		e.ShowId = sql.NullInt64{Valid: false}
	} else {
		e.ShowId = sql.NullInt64{Valid: true, Int64: *episode.ShowId}
	}

	if episode.Title == nil {
		e.Title = sql.NullString{Valid: false}
	} else {
		e.Title = sql.NullString{Valid: true, String: *episode.Title}
	}

	if episode.Number == nil {
		e.Number = sql.NullInt64{Valid: false}
	} else {
		e.Number = sql.NullInt64{Valid: true, Int64: *episode.Number}
	}

	if episode.Season == nil {
		e.Season = sql.NullInt64{Valid: false}
	} else {
		e.Season = sql.NullInt64{Valid: true, Int64: *episode.Season}
	}

	if episode.Length == nil {
		e.Length = sql.NullFloat64{Valid: false}
	} else {
		e.Length = sql.NullFloat64{Valid: true, Float64: *episode.Length}
	}

	return nil
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
