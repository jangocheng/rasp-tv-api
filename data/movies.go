package data

import (
	"database/sql"
	"encoding/json"
)

// Movie represents a movie record from the database
type Movie struct {
	ID        int64
	Title     sql.NullString
	Filepath  string
	Length    sql.NullFloat64
	IsIndexed bool
}

// Methods below are used so we don't send sql.Null* values back

// MarshalJSON implements Marshaller interface
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
		ID        int64    `json:"id"`
		Title     *string  `json:"title"`
		Filepath  string   `json:"filepath"`
		Length    *float64 `json:"length"`
		IsIndexed bool     `json:"isIndexed"`
	}{
		m.ID,
		title,
		m.Filepath,
		length,
		m.IsIndexed,
	}

	return json.Marshal(&mov)
}

// UnmarshalJSON implements Unmarshaller interface
func (m *Movie) UnmarshalJSON(data []byte) error {
	var mov struct {
		ID        int64
		Title     *string
		Filepath  string
		Length    *float64
		IsIndexed bool
	}

	if err := json.Unmarshal(data, &mov); err != nil {
		return err
	}

	m.ID = mov.ID
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
