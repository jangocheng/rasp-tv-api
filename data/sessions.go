package data

import (
	"database/sql"
	"encoding/json"
)

type Session struct {
	Id        int64
	MovieId   sql.NullInt64
	EpisodeId sql.NullInt64
	IsPaused  bool
	IsPlaying bool
	Pid       sql.NullInt64
}

func (s *Session) MarshalJSON() ([]byte, error) {
	var movieId *int64
	if s.MovieId.Valid {
		movieId = &s.MovieId.Int64
	} else {
		movieId = nil
	}

	var episodeId *int64
	if s.EpisodeId.Valid {
		episodeId = &s.EpisodeId.Int64
	} else {
		episodeId = nil
	}

	var pid *int64
	if s.Pid.Valid {
		pid = &s.Pid.Int64
	} else {
		pid = nil
	}

	session := struct {
		Id        int64  `json:"id"`
		MovieId   *int64 `json:"movieId"`
		EpisodeId *int64 `json:"episodeId"`
		IsPaused  bool   `json:"isPaused"`
		IsPlaying bool   `json:"isPlaying"`
		Pid       *int64 `json:"pid"`
	}{
		s.Id,
		movieId,
		episodeId,
		s.IsPaused,
		s.IsPlaying,
		pid,
	}

	return json.Marshal(&session)
}

func (s *Session) UnmarshalJSON(data []byte) error {
	var session struct {
		Id        int64
		MovieId   *int64
		EpisodeId *int64
		IsPaused  bool
		IsPlaying bool
		Pid       *int64
	}

	if err := json.Unmarshal(data, &session); err != nil {
		return err
	}

	s.Id = session.Id
	s.IsPaused = session.IsPaused
	s.IsPlaying = session.IsPlaying

	if session.MovieId == nil {
		s.MovieId = sql.NullInt64{Valid: false}
	} else {
		s.MovieId = sql.NullInt64{Valid: true, Int64: *session.MovieId}
	}

	if session.EpisodeId == nil {
		s.EpisodeId = sql.NullInt64{Valid: false}
	} else {
		s.EpisodeId = sql.NullInt64{Valid: true, Int64: *session.EpisodeId}
	}

	if session.Pid == nil {
		s.Pid = sql.NullInt64{Valid: false}
	} else {
		s.Pid = sql.NullInt64{Valid: true, Int64: *session.Pid}
	}

	return nil
}

func (s *Session) Save(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO session (movieId, episodeId, isPaused, isPlaying) VALUES (?, ?, ?, ?)", s.MovieId, s.EpisodeId, s.IsPaused, s.IsPlaying)
	return err
}

func GetSession(db *sql.DB) (*Session, error) {
	session := Session{}
	err := db.QueryRow("SELECT id, movieId, episodeId, isPaused, isPlaying FROM session ORDER BY id DESC LIMIT 1").Scan(&session.Id, &session.MovieId, &session.EpisodeId, &session.IsPaused, &session.IsPlaying)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &session, nil
}

func ClearSessions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM session")
	return err
}
