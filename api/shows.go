package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"simongeeks.com/joe/rasp-tv/data"

	"github.com/gorilla/mux"
)

func GetAllEpisodes(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var episodes []data.Episode
	var err error

	isIndexedParam := req.URL.Query().Get("isIndexed")
	if len(isIndexedParam) != 0 {
		if isIndexedParam == "true" {
			episodes, err = data.GetEpisodes("WHERE isIndexed = 1 ORDER BY title", context.Db)
		} else {
			episodes, err = data.GetEpisodes("WHERE isIndexed = 0", context.Db)
		}
	} else {
		episodes, err = data.GetEpisodes("", context.Db)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, episodes, nil
}

func GetEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]
	episodes, err := data.GetEpisodes("WHERE id = "+id, context.Db)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(episodes) != 1 {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %s", id)
	}

	return http.StatusOK, &episodes[0], nil
}

func PlayEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]
	episodes, err := data.GetEpisodes("WHERE id = "+id, context.Db)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(episodes) != 1 {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %s", id)
	}

	pid, err := startPlayer(episodes[0].Filepath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	session := data.Session{
		EpisodeId: sql.NullInt64{Int64: episodes[0].Id, Valid: true},
		IsPlaying: true,
		IsPaused:  false,
		Pid:       sql.NullInt64{Int64: pid, Valid: true},
	}
	if err = session.Save(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("Playing episode at %s", episodes[0].Filepath)), nil
}

// func StreamEpisode(r render.Render, params martini.Params, res http.ResponseWriter, req *http.Request, db *sql.DB, logger *log.Logger) {
// 	episodes, err := data.GetEpisodes("WHERE id = "+params["id"], db)
//
// 	if err != nil {
// 		logger.Println(errorMsg(err.Error()))
// 		r.JSON(500, errorResponse(err))
// 		return
// 	}
//
// 	if len(episodes) != 1 {
// 		err = fmt.Errorf("Could not find episode with id: %s", params["id"])
// 		logger.Println(errorMsg(err.Error()))
// 		r.JSON(404, errorResponse(err))
// 		return
// 	}
//
// 	http.ServeFile(res, req, episodes[0].Filepath)
// }

func DeleteEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]

	var err error
	if _, err = strconv.ParseInt(id, 10, 64); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	episodes, err := data.GetEpisodes("WHERE id = "+id, context.Db)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(episodes) != 1 {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %s", id)
	}

	if deleteFile {
		if err = os.Remove(episodes[0].Filepath); err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	if err = episodes[0].DeleteEpisode(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Deleted episode"), nil
}

func AddShow(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	show := &data.Show{}
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(show); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	var err error
	if _, err = show.Add(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, show, nil
}

func GetShows(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	shows, err := data.GetShows("ORDER BY title", context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	allParam := req.URL.Query().Get("all")
	if len(allParam) != 0 && allParam == "true" {
		for i, show := range shows {
			episodes, err := data.GetEpisodes(fmt.Sprintf("WHERE showId = %d", show.Id), context.Db)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}

			shows[i].Episodes = episodes
		}
	}

	return http.StatusOK, shows, nil
}

func GetShow(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id := mux.Vars(req)["id"]
	shows, err := data.GetShows("WHERE id = "+id, context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(shows) != 1 {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find show with id: %s", id)
	}

	show := shows[0]
	episodes, err := data.GetEpisodes("WHERE showId = "+id, context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	show.Episodes = episodes
	return http.StatusOK, show, nil
}

func SaveEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	episode := &data.Episode{}

	if err := json.NewDecoder(req.Body).Decode(episode); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := episode.Update(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("%s saved successfully", episode.Title.String)), nil
}

func StreamEpisode(context *Context, rw http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	episodes, err := data.GetEpisodes("WHERE id = "+id, context.Db)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	if len(episodes) != 1 {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - Could not find episode with id: %s\n", id)
		return
	}

	http.ServeFile(rw, req, episodes[0].Filepath)
}
