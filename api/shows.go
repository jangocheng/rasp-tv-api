package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"simongeeks.com/joe/rasp-tv/data"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func GetAllEpisodes(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	var episodes []data.Episode
	var err error

	isIndexedParam := req.URL.Query().Get("isIndexed")
	if len(isIndexedParam) != 0 {
		if isIndexedParam == "true" {
			episodes, err = data.GetEpisodes("WHERE isIndexed = 1 ORDER BY title", db)
		} else {
			episodes, err = data.GetEpisodes("WHERE isIndexed = 0", db)
		}
	} else {
		episodes, err = data.GetEpisodes("", db)
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, episodes)
}

func GetEpisode(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	episodes, err := data.GetEpisodes("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(episodes) != 1 {
		err = fmt.Errorf("Could not find episode with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	r.JSON(200, &episodes[0])
}

func PlayEpisode(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	episodes, err := data.GetEpisodes("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(episodes) != 1 {
		err = fmt.Errorf("Could not find episode with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	pid, err := startPlayer(episodes[0].Filepath)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	session := data.Session{
		EpisodeId: sql.NullInt64{Int64: episodes[0].Id, Valid: true},
		IsPlaying: true,
		IsPaused:  false,
		Pid:       sql.NullInt64{Int64: pid, Valid: true},
	}
	if err = session.Save(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse(fmt.Sprintf("Playing episode at %s", episodes[0].Filepath)))
}

func StreamEpisode(r render.Render, params martini.Params, res http.ResponseWriter, req *http.Request, db *sql.DB, logger *log.Logger) {
	episodes, err := data.GetEpisodes("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(episodes) != 1 {
		err = fmt.Errorf("Could not find episode with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	http.ServeFile(res, req, episodes[0].Filepath)
}

func DeleteEpisode(r render.Render, req *http.Request, params martini.Params, db *sql.DB, logger *log.Logger) {
	var err error
	if _, err = strconv.ParseInt(params["id"], 10, 64); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	episodes, err := data.GetEpisodes("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(episodes) != 1 {
		err = fmt.Errorf("Could not find episode with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	if deleteFile {
		if err = os.Remove(episodes[0].Filepath); err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, errorResponse(err))
			return
		}
	}

	if err = episodes[0].DeleteEpisode(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse("Deleted episode"))
}

func AddShow(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	show := &data.Show{}
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(show); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	var err error
	if _, err = show.Add(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, show)
}

func GetShows(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	shows, err := data.GetShows("ORDER BY title", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	allParam := req.URL.Query().Get("all")
	if len(allParam) != 0 && allParam == "true" {
		for i, show := range shows {
			episodes, err := data.GetEpisodes(fmt.Sprintf("WHERE showId = %d", show.Id), db)
			if err != nil {
				logger.Println(errorMsg(err.Error()))
				r.JSON(500, errorResponse(err))
				return
			}

			shows[i].Episodes = episodes
		}
	}

	r.JSON(200, shows)
}

func GetShow(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	shows, err := data.GetShows("WHERE id = "+params["id"], db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if len(shows) != 1 {
		err = fmt.Errorf("Could not find show with id: %s", params["id"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(404, errorResponse(err))
		return
	}

	show := shows[0]
	episodes, err := data.GetEpisodes("WHERE showId = "+params["id"], db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	show.Episodes = episodes
	r.JSON(200, &show)
}

func SaveEpisode(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	episode := &data.Episode{}

	if err := json.NewDecoder(req.Body).Decode(episode); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if err := episode.Update(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse(fmt.Sprintf("%s saved successfully", episode.Title.String)))
}
