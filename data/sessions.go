package data

import (
	"database/sql"
	"encoding/json"
)

// Session represents a session record from the database
type Session struct {
	ID        int64
	MovieID   sql.NullInt64
	EpisodeID sql.NullInt64
	IsPaused  bool
	IsPlaying bool
	Pid       sql.NullInt64
}

// Methods below are used so we don't send sql.Null* values back

// MarshalJSON implements Marshaller interface
func (s *Session) MarshalJSON() ([]byte, error) {
	var movieID *int64
	if s.MovieID.Valid {
		movieID = &s.MovieID.Int64
	} else {
		movieID = nil
	}

	var episodeID *int64
	if s.EpisodeID.Valid {
		episodeID = &s.EpisodeID.Int64
	} else {
		episodeID = nil
	}

	var pid *int64
	if s.Pid.Valid {
		pid = &s.Pid.Int64
	} else {
		pid = nil
	}

	session := struct {
		ID        int64  `json:"id"`
		MovieID   *int64 `json:"movieId"`
		EpisodeID *int64 `json:"episodeId"`
		IsPaused  bool   `json:"isPaused"`
		IsPlaying bool   `json:"isPlaying"`
		Pid       *int64 `json:"pid"`
	}{
		s.ID,
		movieID,
		episodeID,
		s.IsPaused,
		s.IsPlaying,
		pid,
	}

	return json.Marshal(&session)
}

// UnmarshalJSON implements Unmarshaller interface
func (s *Session) UnmarshalJSON(data []byte) error {
	var session struct {
		ID        int64
		MovieID   *int64
		EpisodeID *int64
		IsPaused  bool
		IsPlaying bool
		Pid       *int64
	}

	if err := json.Unmarshal(data, &session); err != nil {
		return err
	}

	s.ID = session.ID
	s.IsPaused = session.IsPaused
	s.IsPlaying = session.IsPlaying

	if session.MovieID == nil {
		s.MovieID = sql.NullInt64{Valid: false}
	} else {
		s.MovieID = sql.NullInt64{Valid: true, Int64: *session.MovieID}
	}

	if session.EpisodeID == nil {
		s.EpisodeID = sql.NullInt64{Valid: false}
	} else {
		s.EpisodeID = sql.NullInt64{Valid: true, Int64: *session.EpisodeID}
	}

	if session.Pid == nil {
		s.Pid = sql.NullInt64{Valid: false}
	} else {
		s.Pid = sql.NullInt64{Valid: true, Int64: *session.Pid}
	}

	return nil
}
