package api

import (
	"net/http"
	"sort"

	"simongeeks.com/joe/rasp-tv/data"
)

func ScanMovies(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	moviePaths, err := findVideoFiles(context.Config.MoviePath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	movies, err := data.GetMovies("ORDER BY filepath", context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	for _, path := range moviePaths {
		if index := findMovieByFilePath(movies, len(movies), path); index == -1 {
			_, err := context.Db.Exec("INSERT INTO movies (filepath, isIndexed) VALUES (?, 0)", path)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	return http.StatusOK, statusResponse("Success"), nil
}

func ScanEpisodes(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	showsPaths, err := findVideoFiles(context.Config.ShowsPath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	episodes, err := data.GetEpisodes("ORDER BY filepath", context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	for _, path := range showsPaths {
		if index := findEpisodeByFilePath(episodes, len(episodes), path); index == -1 {
			_, err := context.Db.Exec("INSERT INTO episodes (filepath, isIndexed) VALUES (?, 0)", path)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	return http.StatusOK, statusResponse("Success"), nil
}

func findMovieByFilePath(movies []data.Movie, numMovies int, path string) int {
	index := sort.Search(numMovies, func(i int) bool {
		return movies[i].Filepath >= path
	})

	if index < numMovies && movies[index].Filepath == path {
		return index
	}

	return -1
}

func findEpisodeByFilePath(episodes []data.Episode, numEpisodes int, path string) int {
	index := sort.Search(numEpisodes, func(i int) bool {
		return episodes[i].Filepath >= path
	})

	if index < numEpisodes && episodes[index].Filepath == path {
		return index
	}

	return -1
}
