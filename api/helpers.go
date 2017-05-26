package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// isVideoFile determines if file is video based on extension
func isVideoFile(filename string) bool {
	supportedTypes := []string{"mp4", "mkv", "avi", "m4v", "mov"}
	ext := filepath.Ext(filename)[1:]
	for _, t := range supportedTypes {
		if ext == t {
			return true
		}
	}

	return false
}

func parseIDFromReq(req *http.Request) (int64, error) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// findVideoFiles walks the directory path and finds a list of video files
func findVideoFiles(rootPath string) ([]string, error) {
	videos := make([]string, 0, 70)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !strings.HasPrefix(path, ".") && isVideoFile(path) {
			videos = append(videos, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return videos, nil
}

// statusResponse helper for creating JSON status responses
func statusResponse(msg string) map[string]string {
	return map[string]string{"status": msg}
}
