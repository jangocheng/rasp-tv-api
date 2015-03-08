package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"simongeeks.com/joe/rasp-tv/data"
)

func GetAllMovies(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	isIndexed := req.URL.Query().Get("isIndexed") == "true"
	var movies []data.Movie
	var err error

	if isIndexed {
		movies, err = data.GetMovies("WHERE isIndexed = 1 ORDER BY title", db)
	} else {
		movies, err = data.GetMovies("WHERE isIndexed = 0", db)
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, movies)
}

func SaveMovie(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	movie := data.Movie{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&movie); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if err := movie.Update(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse(fmt.Sprintf("%s saved successfully", movie.Title.String)))
}

func GetMovie(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	movies, err := data.GetMovies("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	r.JSON(200, movies[0])
}

func PlayMovie(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	movies, err := data.GetMovies("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	if err = startPlayer(movies[0].Filepath); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	session := data.Session{MovieId: sql.NullInt64{Int64: movies[0].Id, Valid: true}, IsPlaying: true, IsPaused: false}
	if err = session.Save(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse(fmt.Sprintf("Playing movie at %s", movies[0].Filepath)))
}

func StreamMovie(r render.Render, params martini.Params, res http.ResponseWriter, req *http.Request, db *sql.DB, logger *log.Logger) {
	movies, err := data.GetMovies("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	http.ServeFile(res, req, movies[0].Filepath)
}

func DeleteMovie(r render.Render, req *http.Request, params martini.Params, db *sql.DB, logger *log.Logger) {
	var err error
	if _, err = strconv.ParseInt(params["id"], 10, 64); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	movies, err := data.GetMovies("WHERE id = "+params["id"], db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	if deleteFile {
		if err = os.Remove(movies[0].Filepath); err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, errorResponse(err))
			return
		}
	}

	if err = movies[0].Delete(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse("Deleted movie"))
}
