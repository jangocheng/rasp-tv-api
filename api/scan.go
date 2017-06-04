package api

import (
	"net/http"
	"sort"

	"github.com/simonjm/rasp-tv-api/data"
)

// ScanMovies route that walks the movie directory for new video files
func ScanMovies(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	// get list of new video files in the movie directory
	moviePaths, err := findVideoFiles(context.Config.MoviePath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	db := context.Db
	// sort the list by filepath so we can use binary search to determines
	// if the file needs to be indexed or not
	// a set would be better but go doesn't include one in the stdlib
	movies, err := db.GetMovies("ORDER BY filepath")
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// go through each video file in the directory and test if its already in the list
	for _, path := range moviePaths {
		if index := findMovieByFilePath(movies, len(movies), path); index == -1 {
			// add a new movie record for the new file
			m := data.Movie{Filepath: path, IsIndexed: false}
			if err := db.AddMovie(&m); err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	return http.StatusOK, statusResponse("Success"), nil
}

// ScanEpisodes route that walks the shows directory for new video files
func ScanEpisodes(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	// get list of new video files in the shows directory
	showsPaths, err := findVideoFiles(context.Config.ShowsPath)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	db := context.Db
	// sort the list by filepath so we can use binary search to determines
	// if the file needs to be indexed or not
	// a set would be better but go doesn't include one in the stdlib
	episodes, err := db.GetEpisodes("ORDER BY filepath")
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// go through each video file in the directory and test if its already in the list
	for _, path := range showsPaths {
		if index := findEpisodeByFilePath(episodes, len(episodes), path); index == -1 {
			// add a new episode record for the new file
			e := data.Episode{Filepath: path, IsIndexed: false}
			if err := db.AddEpisode(&e); err != nil {
				return http.StatusInternalServerError, nil, err
			}
		}
	}

	return http.StatusOK, statusResponse("Success"), nil
}

// findMovieByFilePath uses binary search to determine if the path is already in the list of movies
func findMovieByFilePath(movies []data.Movie, numMovies int, path string) int {
	index := sort.Search(numMovies, func(i int) bool {
		return movies[i].Filepath >= path
	})

	if index < numMovies && movies[index].Filepath == path {
		return index
	}

	return -1
}

// findMovieByFilePath uses binary search to determine if the path is already in the list of episodes
func findEpisodeByFilePath(episodes []data.Episode, numEpisodes int, path string) int {
	index := sort.Search(numEpisodes, func(i int) bool {
		return episodes[i].Filepath >= path
	})

	if index < numEpisodes && episodes[index].Filepath == path {
		return index
	}

	return -1
}
