package main

import (
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"database/sql"
	"errors"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"log"
	"os"
	"os/user"
	"simongeeks.com/joe/rasp-tv/api"
)

func getConfig() (*api.Config, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	switch user.Name {
	case "Joe Simon": // macbook
		return &api.Config{
			MoviePath:    "/Volumes/My Passport/Movies",
			ShowsPath:    "/Volumes/My Passport/TV Shows",
			IsProduction: martini.Env == "production",
			LogPath:      "logs.txt",
			DbPath:       "raspTv.db",
			Root:         "/Users/Joe/Projects/go/src/simongeeks.com/joe/rasp-tv",
		}, nil
	case "Joe": // raspberry pi
		return &api.Config{
			MoviePath:    "/media/passport/Movies",
			ShowsPath:    "/media/passport/TV Shows",
			IsProduction: martini.Env == "production",
			LogPath:      "/var/log/rasp-tv/logs.txt",
			DbPath:       "/home/joe/data/raspTv.db",
			Root:         "/home/joe/go/src/simongeeks.com/joe/rasp-tv",
		}, nil
	}

	return nil, errors.New("Could not find a username with a config file")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config, err := getConfig()
	check(err)

	logFile, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	check(err)

	defer logFile.Close()
	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", config.DbPath)
	check(err)
	defer db.Close()

	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Static(config.Root+"/assets", martini.StaticOptions{SkipLogging: true}))
	m.Use(render.Renderer(render.Options{Delims: render.Delims{"[[", "]]"}, Directory: config.Root + "/views"}))
	m.Map(db)
	m.Map(logger)
	m.Map(config)

	router := martini.NewRouter()

	router.Get("/", api.Index)
	router.Group("/scan", func(r martini.Router) {
		r.Get("/movies", api.ScanMovies)
		r.Get("/episodes", api.ScanEpisodes)
	})
	router.Get("/movies", api.GetAllMovies)
	router.Group("/movies", func(r martini.Router) {
		r.Get("/:id", api.GetMovie)
		r.Get("/:id/play", api.PlayMovie)
		r.Get("/:id/stream", api.StreamMovie)
		r.Post("/:id", api.SaveMovie)
		r.Delete("/:id", api.DeleteMovie)
	})
	router.Get("/shows", api.GetShows)
	router.Group("/shows", func(r martini.Router) {
		r.Get("/:id", api.GetShow)
		r.Post("/add", api.AddShow)
		r.Group("/episodes", func(episodeRouter martini.Router) {
			episodeRouter.Get("/:id", api.GetEpisode)
			episodeRouter.Get("/:id/play", api.PlayEpisode)
			episodeRouter.Get("/:id/stream", api.StreamEpisode)
			episodeRouter.Post("/:id", api.SaveEpisode)
			episodeRouter.Delete("/:id", api.DeleteEpisode)
		})
	})

	router.Get("/episodes", api.GetAllEpisodes)
	router.Get("/player/:command", api.RunPlayerCommand)
	// router.Get("/auto", api.AutoIndex)

	m.Action(router.Handle)
	m.Run()
}
