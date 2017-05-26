package data

import (
	"database/sql"
	"encoding/json"
)

// Show represents a show record from the database
type Show struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Episodes []Episode `json:"episodes,omitempty"`
}

// Methods below are used so we don't send sql.Null* values back

// Episode represents an episode record from the database
type Episode struct {
	ID        int64
	ShowID    sql.NullInt64
	Title     sql.NullString
	Number    sql.NullInt64
	Season    sql.NullInt64
	Filepath  string
	Length    sql.NullFloat64
	IsIndexed bool
}

// MarshalJSON implements Marshaller interface
func (e *Episode) MarshalJSON() ([]byte, error) {
	var showID *int64
	if e.ShowID.Valid {
		showID = &e.ShowID.Int64
	} else {
		showID = nil
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
		ID        int64    `json:"id"`
		ShowID    *int64   `json:"showId"`
		Title     *string  `json:"title"`
		Number    *int64   `json:"number"`
		Season    *int64   `json:"season"`
		Filepath  string   `json:"filepath"`
		Length    *float64 `json:"length"`
		IsIndexed bool     `json:"isIndexed"`
	}{
		e.ID,
		showID,
		title,
		number,
		season,
		e.Filepath,
		length,
		e.IsIndexed,
	}

	return json.Marshal(&episode)
}

// UnmarshalJSON implements Unmarshaller interface
func (e *Episode) UnmarshalJSON(data []byte) error {
	var episode struct {
		ID        int64
		ShowID    *int64
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

	e.ID = episode.ID
	e.Filepath = episode.Filepath
	e.IsIndexed = episode.IsIndexed

	if episode.ShowID == nil {
		e.ShowID = sql.NullInt64{Valid: false}
	} else {
		e.ShowID = sql.NullInt64{Valid: true, Int64: *episode.ShowID}
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
