package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func sqlEscape(str string) string {
	return strings.Replace(str, "'", "''", -1)
}

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

func errorMsg(msg string) string {
	return fmt.Sprintf("[Error] %s", msg)
}
