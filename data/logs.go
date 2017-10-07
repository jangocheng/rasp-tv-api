package data

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Log represents a log row in the database
type Log struct {
	ID       int64
	Level    string
	LogDate  time.Time
	Message  string
	Metadata sql.NullString
}

// MarshalJSON implements Marshaller interface
func (l *Log) MarshalJSON() ([]byte, error) {
	log := struct {
		ID       int64     `json:"id"`
		Level    string    `json:"level"`
		LogDate  time.Time `json:"logDate"`
		Message  string    `json:"message"`
		Metadata *string   `json:"metadata"`
	}{
		l.ID,
		l.Level,
		l.LogDate,
		l.Message,
		nil,
	}

	if l.Metadata.Valid {
		log.Metadata = &l.Metadata.String
	}

	return json.Marshal(log)
}

// UnmarshalJSON implements Unmarshaller interface
func (l *Log) UnmarshalJSON(data []byte) error {
	var log struct {
		ID       int64     `json:"id"`
		Level    string    `json:"level"`
		LogDate  time.Time `json:"logDate"`
		Message  string    `json:"message"`
		Metadata *string   `json:"metadata,omitempty"`
	}

	if err := json.Unmarshal(data, &log); err != nil {
		return err
	}

	l.ID = log.ID
	l.Level = log.Level
	l.LogDate = log.LogDate
	l.Message = log.Message

	if log.Metadata != nil {
		l.Metadata = sql.NullString{Valid: true, String: *log.Metadata}
	} else {
		l.Metadata = sql.NullString{Valid: false}
	}

	return nil
}
