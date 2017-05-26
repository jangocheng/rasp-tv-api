package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/simonjm/rasp-tv/data"
)

// GetAllMovies route for getting all of the movies from the databse
func GetAllMovies(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var movies []data.Movie
	var err error

	db := context.Db

	isIndexedParam := req.URL.Query().Get("isIndexed")
	if len(isIndexedParam) != 0 {
		if isIndexedParam == "true" {
			movies, err = db.GetMovies("WHERE isIndexed = 1 ORDER BY title")
		} else {
			movies, err = db.GetMovies("WHERE isIndexed = 0")
		}
	} else {
		movies, err = db.GetMovies("")
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, movies, nil
}

// SaveMovie route for saving a movie to the database. The data to save is deserialized from the request body
func SaveMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	movie := &data.Movie{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(movie); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := context.Db.UpdateMovie(movie); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("%s saved successfully", movie.Title.String)), nil
}

// GetMovie route to get a single movie by id
func GetMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	movie, err := context.Db.GetMovieByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if movie == nil {
		err = fmt.Errorf("Could not find movie with id: %d", id)
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, movie, nil
}

// PlayMovie route that plays a movie with omxplayer
func PlayMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	movie, err := context.Db.GetMovieByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if movie == nil {
		err = fmt.Errorf("Could not find movie with id: %d", id)
		return http.StatusNotFound, nil, err
	}

	pid, err := startPlayer(movie.Filepath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	session := data.Session{
		MovieID:   sql.NullInt64{Int64: movie.ID, Valid: true},
		IsPlaying: true,
		IsPaused:  false,
		Pid:       sql.NullInt64{Int64: pid, Valid: true},
	}

	if err = context.Db.SaveSession(&session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("Playing movie at %s", movie.Filepath)), nil
}

// DeleteMovie deletes a movie by id
func DeleteMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	movie, err := context.Db.GetMovieByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if movie == nil {
		err = fmt.Errorf("Could not find movie with id: %d", id)
		return http.StatusNotFound, nil, err
	}

	if deleteFile {
		if err = os.Remove(movie.Filepath); err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	if err = context.Db.DeleteMovie(movie); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Deleted movie"), nil
}

// StreamMovie route that streams a movie file
func StreamMovie(context *Context, rw http.ResponseWriter, req *http.Request) {
	id, err := parseIDFromReq(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	movie, err := context.Db.GetMovieByID(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	if movie == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - Could not find movie with id: %d\n", id)
		return
	}

	http.ServeFile(rw, req, movie.Filepath)
}
