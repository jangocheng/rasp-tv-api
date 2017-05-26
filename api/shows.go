package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/simonjm/rasp-tv/data"
)

// GetAllEpisodes route for gettting all episodes from the database
func GetAllEpisodes(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var episodes []data.Episode
	var err error

	db := context.Db
	isIndexedParam := req.URL.Query().Get("isIndexed")
	if len(isIndexedParam) != 0 {
		if isIndexedParam == "true" {
			episodes, err = db.GetEpisodes("WHERE isIndexed = 1 ORDER BY title")
		} else {
			episodes, err = db.GetEpisodes("WHERE isIndexed = 0")
		}
	} else {
		episodes, err = db.GetEpisodes("")
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, episodes, nil
}

// GetEpisode route for getting a single episode by id
func GetEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	episode, err := context.Db.GetEpisodeByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if episode == nil {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %d", id)
	}

	return http.StatusOK, episode, nil
}

// PlayEpisode route that plays the episode with omxplayer
func PlayEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	episode, err := context.Db.GetEpisodeByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if episode == nil {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %d", id)
	}

	pid, err := startPlayer(episode.Filepath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	session := data.Session{
		EpisodeID: sql.NullInt64{Int64: episode.ID, Valid: true},
		IsPlaying: true,
		IsPaused:  false,
		Pid:       sql.NullInt64{Int64: pid, Valid: true},
	}
	if err = context.Db.SaveSession(&session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("Playing episode at %s", episode.Filepath)), nil
}

// DeleteEpisode route that deletes an episode by id
func DeleteEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	episode, err := context.Db.GetEpisodeByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if episode == nil {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find episode with id: %d", id)
	}

	if deleteFile {
		if err = os.Remove(episode.Filepath); err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	if err = context.Db.DeleteEpisode(episode); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Deleted episode"), nil
}

// AddShow route for saving a show to the database. The data to save is deserialized from the request body
func AddShow(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	show := &data.Show{}
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(show); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := context.Db.AddShow(show); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, show, nil
}

// GetShows route to get all of the shows from the database
func GetShows(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	shows, err := context.Db.GetShows("ORDER BY title")
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	allParam := req.URL.Query().Get("all")
	if len(allParam) != 0 && allParam == "true" {
		// if the all parameter is present then return all of the episodes for each show
		for i, show := range shows {
			episodes, err := context.Db.GetEpisodes("WHERE showId = ?", show.ID)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}

			shows[i].Episodes = episodes
		}
	}

	return http.StatusOK, shows, nil
}

// GetShow gets a single show with all of the episodes from the database
func GetShow(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	id, err := parseIDFromReq(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	show, err := context.Db.GetShowByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if show == nil {
		return http.StatusNotFound, nil, fmt.Errorf("Could not find show with id: %d", id)
	}

	episodes, err := context.Db.GetEpisodes("WHERE showId = ?", id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	show.Episodes = episodes
	return http.StatusOK, show, nil
}

// SaveEpisode route for saving an episode to the database. The data to save is deserialized from the request body
func SaveEpisode(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	episode := &data.Episode{}

	if err := json.NewDecoder(req.Body).Decode(episode); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := context.Db.UpdateEpisode(episode); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse(fmt.Sprintf("%s saved successfully", episode.Title.String)), nil
}

// StreamEpisode route that streams the video file for an episode
func StreamEpisode(context *Context, rw http.ResponseWriter, req *http.Request) {
	id, err := parseIDFromReq(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	episode, err := context.Db.GetEpisodeByID(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		context.Logger.Printf("[Error] - %s\n", err)
		return
	}

	if episode == nil {
		rw.WriteHeader(http.StatusNotFound)
		context.Logger.Printf("[Error] - Could not find episode with id: %d\n", id)
		return
	}

	http.ServeFile(rw, req, episode.Filepath)
}
