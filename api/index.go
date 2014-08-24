package api

import (
	"github.com/martini-contrib/render"
)

type Config struct {
	MoviePath    string
	ShowsPath    string
	IsProduction bool
	LogPath      string
	DbPath       string
}

func Index(r render.Render, config *Config) {
	data := struct {
		IsProduction bool
	}{
		config.IsProduction,
	}
	r.HTML(200, "index", data)
}
