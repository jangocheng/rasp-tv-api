package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/simonjm/rasp-tv/api"
	"github.com/simonjm/rasp-tv/data"
)

// getConfig create struct to hold all of the configuration based on current user
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
		}, nil
	default: // raspberry pi
		return &api.Config{
			MoviePath:    "/media/passport/Movies",
			ShowsPath:    "/media/passport/TV Shows",
			IsProduction: os.Getenv("RASPTV_ENV") == "production",
			LogPath:      "/var/log/rasp-tv/logs.txt",
			DbPath:       "/home/joe/data/raspTv.db",
		}, nil
	}
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
	db, err := data.NewRaspTvDatabase(config.DbPath)
	check(err)

	if err := db.ClearSession(); err != nil {
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

	// set up api routes
	router := mux.NewRouter()
	router.Handle("/scan/movies", createHandler(api.ScanMovies)).Methods("GET")
	router.Handle("/scan/episodes", createHandler(api.ScanEpisodes)).Methods("GET")
	router.Handle("/movies", createHandler(api.GetAllMovies)).Methods("GET")
	router.Handle("/movies/{id}", createHandler(api.GetMovie)).Methods("GET")
	router.Handle("/movies/{id}/play", createHandler(api.PlayMovie)).Methods("GET")
	router.Handle("/movies/{id}/stream", createStreamHandler(api.StreamMovie)).Methods("GET")
	router.Handle("/movies/{id}", createHandler(api.SaveMovie)).Methods("POST", "OPTIONS")
	router.Handle("/movies/{id}", createHandler(api.DeleteMovie)).Methods("DELETE", "OPTIONS")
	router.Handle("/shows", createHandler(api.GetShows)).Methods("GET")
	router.Handle("/shows/{id}", createHandler(api.GetShow)).Methods("GET")
	router.Handle("/shows/add", createHandler(api.AddShow)).Methods("POST", "OPTIONS")
	router.Handle("/shows/episodes/{id}", createHandler(api.GetEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}/play", createHandler(api.PlayEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}/stream", createStreamHandler(api.StreamEpisode)).Methods("GET")
	router.Handle("/shows/episodes/{id}", createHandler(api.SaveEpisode)).Methods("POST", "OPTIONS")
	router.Handle("/shows/episodes/{id}", createHandler(api.DeleteEpisode)).Methods("DELETE", "OPTIONS")
	router.Handle("/episodes", createHandler(api.GetAllEpisodes)).Methods("GET")
	router.Handle("/player/command/{command}", createHandler(api.RunPlayerCommand)).Methods("GET")
	router.Handle("/player/session", createHandler(api.NowPlaying)).Methods("GET")
	router.Handle("/player/session", createHandler(api.ClearSession)).Methods("DELETE", "OPTIONS")
	router.Handle("/player/session", createHandler(api.UpdateSession)).Methods("POST", "OPTIONS")

	// get port from the environment
	port := os.Getenv("RASPTV_PORT")
	if len(port) == 0 {
		port = "3000"
	}
	logger.Printf("Server started on port %s\n", port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

// isPlayingPoller clears the session when the omxplayer process that matches
// the pid is not longer running
func isPlayingPoller(dbPath string, logger *log.Logger) {
	db, err := data.NewRaspTvDatabase(dbPath)
	if err != nil {
		logger.Fatal(err)
	}

	for _ = range time.Tick(time.Second * 2) {
		session, err := db.GetSession()
		if err != nil {
			logger.Println(err)
			continue
		}

		if session.Pid.Valid {
			if !data.IsProcessRunning(session.Pid.Int64) {
				logger.Println("omxplayer has exited. Clearing session")
				if err := db.ClearSession(); err != nil {
					logger.Println(err)
				}
			}
		}
	}
}
