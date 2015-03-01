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

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func GetAllEpisodes(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	isIndexed := req.URL.Query().Get("isIndexed") == "true"
	var episodes []data.Episode
	var err error

	if isIndexed {
		episodes, err = data.GetEpisodes("WHERE isIndexed = 1 ORDER BY title", db)
	} else {
		episodes, err = data.GetEpisodes("WHERE isIndexed = 0", db)
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

	r.JSON(200, episodes[0])
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

	if err = startPlayer(episodes[0].Filepath); err != nil {
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
	show := data.Show{}
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(&show); err != nil {
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

	r.JSON(200, statusResponse("Success"))
}

func GetShows(r render.Render, db *sql.DB, logger *log.Logger) {
	shows, err := data.GetShows("", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
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
	r.JSON(200, show)
}

func SaveEpisode(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	episode := data.Episode{}

	if err := json.NewDecoder(req.Body).Decode(&episode); err != nil {
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
