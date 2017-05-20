package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simonjm/rasp-tv/data"
)

func GetAllMovies(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var movies []data.Movie
	var err error

	isIndexedParam := req.URL.Query().Get("isIndexed")
	if len(isIndexedParam) != 0 {
		if isIndexedParam == "true" {
			movies, err = data.GetMovies("WHERE isIndexed = 1 ORDER BY title", context.Db)
		} else {
			movies, err = data.GetMovies("WHERE isIndexed = 0", context.Db)
		}
	} else {
		movies, err = data.GetMovies("", context.Db)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, movies, nil
}

func SaveMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	movie := &data.Movie{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(movie); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := movie.Update(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("%s saved successfully", movie.Title.String)), nil
}

func GetMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]
	movies, err := data.GetMovies("WHERE id = "+id, context.Db)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", id)
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, &movies[0], nil
}

func PlayMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]
	movies, err := data.GetMovies("WHERE id = "+id, context.Db)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", id)
		return http.StatusNotFound, nil, err
	}

	pid, err := startPlayer(movies[0].Filepath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	session := data.Session{
		MovieId:   sql.NullInt64{Int64: movies[0].Id, Valid: true},
		IsPlaying: true,
		IsPaused:  false,
		Pid:       sql.NullInt64{Int64: pid, Valid: true},
	}
	if err = session.Save(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("Playing movie at %s", movies[0].Filepath)), nil
}

func DeleteMovie(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]

	var err error
	if _, err = strconv.ParseInt(id, 10, 64); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	movies, err := data.GetMovies("WHERE id = "+id, context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(movies) != 1 {
		err = fmt.Errorf("Could not find movie with id: %s", id)
		return http.StatusNotFound, nil, err
	}

	if deleteFile {
		if err = os.Remove(movies[0].Filepath); err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	if err = movies[0].Delete(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Deleted movie"), nil
}

func StreamMovie(context *Context, rw http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	movies, err := data.GetMovies("WHERE id = "+id, context.Db)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	if len(movies) != 1 {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - Could not find movie with id: %s\n", id)
		return
	}

	http.ServeFile(rw, req, movies[0].Filepath)
}
