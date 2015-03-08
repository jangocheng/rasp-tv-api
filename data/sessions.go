package data

import "database/sql"

type Session struct {
	Id        int64
	MovieId   sql.NullInt64
	EpisodeId sql.NullInt64
	IsPaused  bool
	IsPlaying bool
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
