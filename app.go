package main

import (
	"database/sql"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mxk/go-sqlite/sqlite3"
	"github.com/rs/cors"
	"simongeeks.com/joe/rasp-tv/api"
	"simongeeks.com/joe/rasp-tv/data"
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
			IsProduction: os.Getenv("RASPTV_ENV") == "production",
			LogPath:      "logs.txt",
			DbPath:       "raspTv.db",
			Root:         "/Users/Joe/Projects/go/src/simongeeks.com/joe/rasp-tv",
		}, nil
	default: // raspberry pi
		return &api.Config{
			MoviePath:    "/media/passport/Movies",
			ShowsPath:    "/media/passport/TV Shows",
			IsProduction: os.Getenv("RASPTV_ENV") == "production",
			LogPath:      "/var/log/rasp-tv/logs.txt",
			DbPath:       "/home/joe/data/raspTv.db",
			Root:         "/home/joe/workspace/go/src/simongeeks.com/joe/rasp-tv",
		}, nil
	}

	// return nil, errors.New("Could not find a username with a config file")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config, err := getConfig()
	check(err)

	// set up logging
	logFile, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	check(err)
	defer logFile.Close()

	var logWriter io.Writer
	if !config.IsProduction {
		logWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		logWriter = logFile
	}
	logger := log.New(logWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	// clear session when app first starts
	db, err := sql.Open("sqlite3", config.DbPath)
	check(err)

	if err := data.ClearSessions(db); err != nil {
		logger.Fatal(err)
	}
	db.Close()

	if config.IsProduction {
		go isPlayingPoller(config.DbPath, logger)
	}

	// handler creators
	createHandler := func(handler api.ApiHandlerFunc) http.Handler {
		return cors.Default().Handler(api.NewApiHandler(logger, config, handler))
	}

	createStreamHandler := func(handler api.StreamHandlerFunc) http.Handler {
		return cors.Default().Handler(api.NewStreamHandler(logger, config, handler))
	}

	staticDir := filepath.Join(config.Root, "static")

	// set up api routes
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		templatePath := filepath.Join(config.Root, "views", "index.tmpl")
		data, err := ioutil.ReadFile(templatePath)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.New("template").Delims("[[", "]]").Parse(string(data))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.Execute(rw, config); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	})
	router.Handle("/scan/movies", createHandler(api.ScanMovies)).Methods("GET")
	router.Handle("/scan/episodes", createHandler(api.ScanEpisodes)).Methods("GET")
	router.Handle("/movies", createHandler(api.GetAllMovies)).Methods("GET")
	router.Handle("/movies/{id}", createHandler(api.GetMovie)).Methods("GET")
	router.Handle("/movies/{id}/play", createHandler(api.PlayMovie)).Methods("GET")
	router.Handle("/movies/{id}/stream", createStreamHandler(api.StreamMovie)).Methods("GET")
	router.Handle("/movies/{id}", createHandler(api.SaveMovie)).Methods("POST")
	router.Handle("/movies/{id}", createHandler(api.DeleteMovie)).Methods("DELETE")
	router.Handle("/shows", createHandler(api.GetShows)).Methods("GET")
	router.Handle("/shows/{id}", createHandler(api.GetShow)).Methods("GET")
	router.Handle("/shows/add", createHandler(api.AddShow)).Methods("POST")
	router.Handle("/shows/episodes/{id}", createHandler(api.GetEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}/play", createHandler(api.PlayEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}/stream", createStreamHandler(api.StreamEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}", createHandler(api.SaveEpisode)).Methods("POST")
	router.Handle("/shows/episodes/{id}", createHandler(api.DeleteEpisode)).Methods("DELETE")
	router.Handle("/episodes", createHandler(api.GetAllEpisodes)).Methods("GET")
	router.Handle("/player/command/{command}", createHandler(api.RunPlayerCommand)).Methods("GET")
	router.Handle("/player/session", createHandler(api.NowPlaying)).Methods("GET")
	router.Handle("/player/session", createHandler(api.ClearSession)).Methods("DELETE")
	router.Handle("/player/session", createHandler(api.UpdateSession)).Methods("POST")

	logger.Fatal(http.ListenAndServe(":3000", router))
}

func isPlayingPoller(dbPath string, logger *log.Logger) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal(err)
	}

	for _ = range time.Tick(time.Second * 2) {
		session, err := data.GetSession(db)
		if err != nil {
			logger.Println(err)
			continue
		}

		if session.Pid.Valid {
			if !data.IsProcessRunning(session.Pid.Int64) {
				logger.Println("omxplayer has exited. Clearing session")
				if err := data.ClearSessions(db); err != nil {
					logger.Println(err)
				}
			}
		}
	}
}
