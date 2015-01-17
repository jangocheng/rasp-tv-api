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
	"os"
	"strconv"
)

type Show struct {
	Id       int64
	Title    string
	Episodes []Episode
}

func (s *Show) Add(db *sql.DB) (int64, error) {
	query := fmt.Sprintf("INSERT INTO shows (title) VALUES ('%s')", sqlEscape(s.Title))
	result, err := db.Exec(query)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, err
}

type Episode struct {
	Id        int64
	ShowId    sql.NullInt64
	Title     sql.NullString
	Number    sql.NullInt64
	Season    sql.NullInt64
	Filepath  string
	Length    sql.NullFloat64
	IsIndexed bool
}

func (e *Episode) Update(db *sql.DB) error {
	if !e.ShowId.Valid {
		return errors.New("Cannot update episode with invalid showId")
	}

	if !e.Title.Valid {
		return errors.New("Cannot update episode with invalid title")
	}

	if !e.Number.Valid {
		return errors.New("Cannot update episode with invalid episode number")
	}

	if !e.Season.Valid {
		return errors.New("Cannot update episode with invalid season")
	}

	query := fmt.Sprintf("UPDATE episodes SET showId = %d, title = '%s', episodeNumber = %d, season = %d, isIndexed = 1 WHERE id = %d", e.ShowId.Int64, sqlEscape(e.Title.String), e.Number.Int64, e.Season.Int64, e.Id)
	_, err := db.Exec(query)
	return err
}

func GetAllEpisodes(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	isIndexed := req.URL.Query().Get("isIndexed") == "true"
	var episodes []Episode
	var err error

	if isIndexed {
		episodes, err = getEpisodesFromDb("WHERE isIndexed = 1 ORDER BY title", db)
	} else {
		episodes, err = getEpisodesFromDb("WHERE isIndexed = 0", db)
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, episodes)
}

func GetEpisode(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	episodes, err := getEpisodesFromDb("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(episodes) != 1 {
		msg := "Could not find episode with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	r.JSON(200, episodes[0])
}

func PlayEpisode(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	episodes, err := getEpisodesFromDb("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(episodes) != 1 {
		msg := "Could not find episode with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	if err = startPlayer(episodes[0].Filepath); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, fmt.Sprintf("Playing episode at %s", episodes[0].Filepath))
}

func StreamEpisode(r render.Render, params martini.Params, res http.ResponseWriter, req *http.Request, db *sql.DB, logger *log.Logger) {
	episodes, err := getEpisodesFromDb("WHERE id = "+params["id"], db)

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(episodes) != 1 {
		msg := "Could not find episode with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	http.ServeFile(res, req, episodes[0].Filepath)
}

func DeleteEpisode(r render.Render, req *http.Request, params martini.Params, db *sql.DB, logger *log.Logger) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	deleteFile := req.URL.Query().Get("file") == "true"
	if deleteFile {
		episodes, err := getEpisodesFromDb("WHERE id = "+params["id"], db)
		if err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}

		if len(episodes) != 1 {
			msg := "Could not find episode with id: " + params["id"]
			logger.Println(errorMsg(msg))
			r.JSON(404, map[string]string{"error": msg})
			return
		}

		if err = os.Remove(episodes[0].Filepath); err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}
	}

	_, err = db.Exec(fmt.Sprintf("DELETE FROM episodes WHERE Id = %d;", id))
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, "Deleted episode")
}

func AddShow(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	show := Show{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&show); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	var err error
	if _, err = show.Add(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, "Success")
}

func GetShows(r render.Render, db *sql.DB, logger *log.Logger) {
	shows, err := getShowsFromDb("", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, shows)
}

func GetShow(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
	shows, err := getShowsFromDb("WHERE id = "+params["id"], db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if len(shows) != 1 {
		msg := "Could not find show with id: " + params["id"]
		logger.Println(errorMsg(msg))
		r.JSON(404, map[string]string{"error": msg})
		return
	}

	show := shows[0]
	episodes, err := getEpisodesFromDb("WHERE showId = "+params["id"], db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	show.Episodes = episodes
	r.JSON(200, show)
}

func SaveEpisode(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	episode := Episode{}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&episode); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if err := episode.Update(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, fmt.Sprintf("%s saved successfully", episode.Title.String))
}

func getShowsFromDb(filter string, db *sql.DB) ([]Show, error) {
	shows := make([]Show, 0, 10)
	rows, err := db.Query("SELECT id, title FROM shows " + filter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return shows, nil
	}

	for rows.Next() {
		show := Show{}
		rows.Scan(&show.Id, &show.Title)
		shows = append(shows, show)
	}

	return shows, nil
}

func getEpisodesFromDb(filter string, db *sql.DB) ([]Episode, error) {
	episodes := make([]Episode, 0, 20)
	rows, err := db.Query("SELECT id, title, episodeNumber, season, filepath, length, isIndexed, showId FROM episodes " + filter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	defer rows.Close()

	if err == io.EOF {
		return episodes, nil
	}

	for rows.Next() {
		e := Episode{}
		if err := rows.Scan(&e.Id, &e.Title, &e.Number, &e.Season, &e.Filepath, &e.Length, &e.IsIndexed, &e.ShowId); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}

	return episodes, nil
}
