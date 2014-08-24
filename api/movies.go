package api

import (
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"net/http"
)

type Movie struct {
	Id        int64
	Title     sql.NullString
	Filepath  string
	Length    float64
	IsIndexed bool
}

func (m *Movie) Update(db *sql.DB) error {
	if !m.Title.Valid {
		return errors.New("Cannot update movie with invalid title")
	}

	query := fmt.Sprintf("UPDATE movies SET title = '%s', isIndexed = 1 WHERE id = %d", sqlEscape(m.Title.String), m.Id)
	_, err := db.Exec(query)
	return err
}

func GetAllMovies(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	isIndexed := req.URL.Query().Get("isIndexed") == "true"
	var movies []Movie
	var err error

	if isIndexed {
		movies, err = getMoviesFromDb("WHERE isIndexed = 1 ORDER BY title", db)
	} else {
		movies, err = getMoviesFromDb("WHERE isIndexed = 0", db)
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, movies)
}

func SaveMovie(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	movie := Movie{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&movie); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if err := movie.Update(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, fmt.Sprintf("%s saved successfully", movie.Title.String))
}

func GetMovie(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	movies, err := getMoviesFromDb("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(movies) != 1 {
		msg := "Could not find movie with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	r.JSON(200, movies[0])
}

func PlayMovie(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	movies, err := getMoviesFromDb("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(movies) != 1 {
		msg := "Could not find movie with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	if err = startPlayer(movies[0].Filepath); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, fmt.Sprintf("Playing movie at %s", movies[0].Filepath))
}

func getMoviesFromDb(filter string, db *sql.DB) ([]Movie, error) {
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
		rows.Scan(&m.Id, &m.Title, &m.Filepath, &m.Length, &m.IsIndexed)
		movies = append(movies, m)
	}

	return movies, nil
}
