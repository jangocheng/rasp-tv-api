package data

import (
	"database/sql"
	"fmt"
	"io"
	"time"
)

// RaspTvDataFetcher interface for interacting with the database
type RaspTvDataFetcher interface {
	GetMovies(filter string, params ...interface{}) ([]Movie, error)
	GetMovieByID(id int64) (*Movie, error)
	AddMovie(movie *Movie) error
	UpdateMovie(movie *Movie) error
	DeleteMovie(movie *Movie) error
	AddShow(show *Show) error
	AddEpisode(episode *Episode) error
	UpdateEpisode(episode *Episode) error
	DeleteEpisode(episode *Episode) error
	GetShows(filter string, params ...interface{}) ([]Show, error)
	GetShowByID(id int64) (*Show, error)
	GetEpisodes(filter string, params ...interface{}) ([]Episode, error)
	GetEpisodeByID(id int64) (*Episode, error)
	SaveSession(session *Session) error
	GetSession() (*Session, error)
	ClearSession() error
	SaveLogs(logs []Log) error
}

// RaspTvDatabase struct with methods for accessing the database
type RaspTvDatabase struct {
	db *sql.DB
}

// NewRaspTvDatabase opens a connection to the sqlite database and constructs the RaspTvDatabase struct
func NewRaspTvDatabase(dbPath string) (*RaspTvDatabase, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &RaspTvDatabase{db: db}, nil
}

// GetMovies gets all movies from the database with the optional filter
func (raspDb *RaspTvDatabase) GetMovies(filter string, params ...interface{}) ([]Movie, error) {
	movies := make([]Movie, 0, 70)
	rows, err := raspDb.db.Query("SELECT id, title, filepath, length, isIndexed FROM movies "+filter, params...)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return movies, nil
	}

	for rows.Next() {
		m := Movie{}
		if err := rows.Scan(&m.ID, &m.Title, &m.Filepath, &m.Length, &m.IsIndexed); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

// GetMovieByID gets a movie from the database by id
func (raspDb *RaspTvDatabase) GetMovieByID(id int64) (*Movie, error) {
	movies, err := raspDb.GetMovies("WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(movies) != 1 {
		return nil, nil
	}

	return &movies[0], nil
}

// AddMovie adds a movie to the database and modifies the ID field with the new inserted ID
func (raspDb *RaspTvDatabase) AddMovie(movie *Movie) error {
	query := "INSERT INTO movies (title, filepath, isIndexed) VALUES (?, ?, ?)"
	result, err := raspDb.db.Exec(query, movie.Title, movie.Filepath, movie.IsIndexed)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	movie.ID = id

	return nil
}

// UpdateMovie updates a movie record in the database
func (raspDb *RaspTvDatabase) UpdateMovie(movie *Movie) error {
	if !movie.Title.Valid {
		return fmt.Errorf("Cannot update movie with invalid title")
	}

	_, err := raspDb.db.Exec("UPDATE movies SET title = ?, isIndexed = 1 WHERE id = ?", movie.Title, movie.ID)
	return err
}

// DeleteMovie deletes a movie record from the database
func (raspDb *RaspTvDatabase) DeleteMovie(movie *Movie) error {
	_, err := raspDb.db.Exec("DELETE FROM movies WHERE Id = ?", movie.ID)
	return err
}

// AddShow adds a show to the database and modifies the ID field with the new inserted ID
func (raspDb *RaspTvDatabase) AddShow(show *Show) error {
	result, err := raspDb.db.Exec("INSERT INTO shows (title) VALUES (?)", show.Title)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	show.ID = id

	return nil
}

// AddEpisode adds an episode to the database and modifies the ID field with the new inserted ID
func (raspDb *RaspTvDatabase) AddEpisode(episode *Episode) error {
	query := `
		INSERT INTO episodes (showId, title, episodeNumber, season, filepath, isIndexed)
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := raspDb.db.Exec(query, episode.ShowID, episode.Title, episode.Number, episode.Season, episode.Filepath, episode.IsIndexed)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	episode.ID = id

	return nil
}

// UpdateEpisode updates an episode record in the database
func (raspDb *RaspTvDatabase) UpdateEpisode(episode *Episode) error {
	if !episode.ShowID.Valid {
		return fmt.Errorf("Cannot update episode with invalid showId")
	}

	if !episode.Title.Valid {
		return fmt.Errorf("Cannot update episode with invalid title")
	}

	if !episode.Number.Valid {
		return fmt.Errorf("Cannot update episode with invalid episode number")
	}

	if !episode.Season.Valid {
		return fmt.Errorf("Cannot update episode with invalid season")
	}

	query := "UPDATE episodes SET showId = ?, title = ?, episodeNumber = ?, season = ?, isIndexed = 1 WHERE id = ?"
	_, err := raspDb.db.Exec(query, episode.ShowID, episode.Title, episode.Number, episode.Season, episode.ID)
	return err
}

// DeleteEpisode deletes an episode from the database
func (raspDb *RaspTvDatabase) DeleteEpisode(episode *Episode) error {
	var err error
	_, err = raspDb.db.Exec("DELETE FROM episodes WHERE Id = ?;", episode.ID)
	return err
}

// GetShows gets all shows from the database with the optional filter
func (raspDb *RaspTvDatabase) GetShows(filter string, params ...interface{}) ([]Show, error) {
	shows := make([]Show, 0, 20)
	rows, err := raspDb.db.Query("SELECT id, title FROM shows "+filter, params...)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return shows, nil
	}

	for rows.Next() {
		show := Show{}
		rows.Scan(&show.ID, &show.Title)
		shows = append(shows, show)
	}

	return shows, nil
}

// GetShowByID gets a show from the database by id
func (raspDb *RaspTvDatabase) GetShowByID(id int64) (*Show, error) {
	shows, err := raspDb.GetShows("WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(shows) != 1 {
		return nil, nil
	}

	return &shows[0], nil
}

// GetEpisodes gets all episodes from the database with the optional filter
func (raspDb *RaspTvDatabase) GetEpisodes(filter string, params ...interface{}) ([]Episode, error) {
	episodes := make([]Episode, 0, 100)
	rows, err := raspDb.db.Query("SELECT id, title, episodeNumber, season, filepath, length, isIndexed, showId FROM episodes "+filter, params...)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return episodes, nil
	}

	for rows.Next() {
		e := Episode{}
		if err := rows.Scan(&e.ID, &e.Title, &e.Number, &e.Season, &e.Filepath, &e.Length, &e.IsIndexed, &e.ShowID); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}

	return episodes, nil
}

// GetEpisodeByID gets an episode from the database by id
func (raspDb *RaspTvDatabase) GetEpisodeByID(id int64) (*Episode, error) {
	episodes, err := raspDb.GetEpisodes("WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(episodes) != 1 {
		return nil, nil
	}

	return &episodes[0], nil
}

// SaveSession inserts a new session record into the database
func (raspDb *RaspTvDatabase) SaveSession(session *Session) error {
	query := "INSERT INTO session (movieId, episodeId, isPaused, isPlaying, pid) VALUES (?, ?, ?, ?, ?)"
	_, err := raspDb.db.Exec(query, session.MovieID, session.EpisodeID, session.IsPaused, session.IsPlaying, session.Pid)
	return err
}

// GetSession gets the session from the database
func (raspDb *RaspTvDatabase) GetSession() (*Session, error) {
	session := Session{}
	query := "SELECT id, movieId, episodeId, isPaused, isPlaying, pid FROM session ORDER BY id DESC LIMIT 1"
	err := raspDb.db.QueryRow(query).Scan(&session.ID, &session.MovieID, &session.EpisodeID, &session.IsPaused, &session.IsPlaying, &session.Pid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &session, nil
}

// ClearSession removes the session from the database
func (raspDb *RaspTvDatabase) ClearSession() error {
	_, err := raspDb.db.Exec("DELETE FROM session")
	return err
}

// SaveLogs inserts a batch of log entries
func (raspDb *RaspTvDatabase) SaveLogs(logs []Log) error {
	tx, err := raspDb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO logs (level, logDate, message, metadata) VALUES (?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, log := range logs {
		date := log.LogDate
		if date.IsZero() {
			date = time.Now()
		}

		var metadata string
		if log.Metadata.Valid {
			metadata = log.Metadata.String
		}

		_, err := stmt.Exec(log.Level, date, log.Message, metadata)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Close closes the inner sqlite connection
func (raspDb *RaspTvDatabase) Close() error {
	return raspDb.db.Close()
}
