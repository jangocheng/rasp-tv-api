package api

import (
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"database/sql"
	"fmt"
	"github.com/martini-contrib/render"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ScanMovies(r render.Render, db *sql.DB, logger *log.Logger, config *Config) {
	moviePaths, err := findVideoFiles(config.MoviePath)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	movies, err := getMoviesFromDb("ORDER BY filepath", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	for _, path := range moviePaths {
		if index := findMovieByFilePath(movies, len(movies), path); index == -1 {
			_, err := db.Exec(fmt.Sprintf("INSERT INTO movies (filepath, isIndexed) VALUES ('%s', 0)", sqlEscape(path)))
			if err != nil {
				logger.Println(errorMsg(err.Error()))
				r.JSON(500, map[string]string{"error": err.Error()})
				return
			}
		}
	}

	r.JSON(200, "Success")
}

func ScanEpisodes(r render.Render, db *sql.DB, logger *log.Logger, config *Config) {
	showsPaths, err := findVideoFiles(config.ShowsPath)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	episodes, err := getEpisodesFromDb("ORDER BY filepath", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	for _, path := range showsPaths {
		if index := findEpisodeByFilePath(episodes, len(episodes), path); index == -1 {
			_, err := db.Exec(fmt.Sprintf("INSERT INTO episodes (filepath, isIndexed) VALUES ('%s', 0)", sqlEscape(path)))
			if err != nil {
				logger.Println(errorMsg(err.Error()))
				r.JSON(500, map[string]string{"error": err.Error()})
				return
			}
		}
	}

	r.JSON(200, "Success")
}

// used to auto index movies and tv shows based on directory structure and filenames
func AutoIndex(r render.Render, db *sql.DB, logger *log.Logger) {
	movies, err := getMoviesFromDb("WHERE isIndexed = 0", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	for _, movie := range movies {
		dir := filepath.Dir(movie.Filepath)
		title := dir[strings.LastIndex(dir, "/")+1:]
		movie.Title.String = title
		movie.Title.Valid = true
		movie.Update(db)
	}

	episodes, err := getEpisodesFromDb("WHERE isIndexed = 0", db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	// /Volumes/pi/TV Shows/Cosmos/1/11 - The Persistence Of Memory.mp4"

	// map of show names to their ids
	showMap := make(map[string]int64)
	shows, err := getShowsFromDb("", db)
	for _, show := range shows {
		showMap[show.Title] = show.Id
	}

	for _, episode := range episodes {
		episodeFile := filepath.Base(episode.Filepath)

		// get episode number
		episodeNumRegex, err := regexp.Compile("^\\d+")

		if err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}

		num, err := strconv.ParseInt(episodeNumRegex.FindString(episodeFile), 10, 64)

		if err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}

		episode.Number.Int64 = num
		episode.Number.Valid = true

		// title
		episodeTitle := episodeNumRegex.ReplaceAllString(episodeFile, "")[2:]
		episode.Title.String = strings.Trim(episodeTitle[0:strings.LastIndex(episodeTitle, ".")], " ")
		episode.Title.Valid = true

		// get season
		seasonRegex, err := regexp.Compile("\\d+$")
		if err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}

		episodeDir := filepath.Dir(episode.Filepath)
		num, err = strconv.ParseInt(seasonRegex.FindString(episodeDir), 10, 64)
		if err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}

		episode.Season.Int64 = num
		episode.Season.Valid = true

		// get show
		episodeDir = seasonRegex.ReplaceAllString(episodeDir, "")
		showName := filepath.Base(episodeDir)

		// create new show if neccesary
		if _, ok := showMap[showName]; !ok {
			show := Show{Title: showName}
			id, err := show.Add(db)
			if err != nil {
				logger.Println(errorMsg(err.Error()))
				r.JSON(500, map[string]string{"error": err.Error()})
				return
			}
			showMap[showName] = id
		}

		episode.ShowId.Int64 = showMap[showName]
		episode.ShowId.Valid = true

		// save episode
		if err = episode.Update(db); err != nil {
			logger.Println(errorMsg(err.Error()))
			r.JSON(500, map[string]string{"error": err.Error()})
			return
		}
	}

	r.JSON(200, "Success")
}

func findMovieByFilePath(movies []Movie, numMovies int, path string) int {
	index := sort.Search(numMovies, func(i int) bool {
		return movies[i].Filepath >= path
	})
	if index < numMovies && movies[index].Filepath == path {
		return index
	} else {
		return -1
	}
}

func findEpisodeByFilePath(episodes []Episode, numEpisodes int, path string) int {
	index := sort.Search(numEpisodes, func(i int) bool {
		return episodes[i].Filepath >= path
	})
	if index < numEpisodes && episodes[index].Filepath == path {
		return index
	} else {
		return -1
	}
}
