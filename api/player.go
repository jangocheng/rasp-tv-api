package api

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"os/exec"
	"strconv"
)

var pipe io.WriteCloser

func startPlayer(path string) error {
	var err error
	if pipe != nil {
		err = stop()
		if err != nil {
			return err
		}
	}

	command := exec.Command("omxplayer", "-o", "hdmi", "-b", path)
	pipe, err = command.StdinPipe()
	if err != nil {
		return err
	}
	return command.Start()
}

func RunPlayerCommand(r render.Render, params martini.Params, logger *log.Logger) {
	if pipe == nil {
		msg := "Player not started"
		logger.Println(errorMsg(msg))
		r.JSON(500, map[string]string{"error": msg})
		return
	}

	cmd, err := strconv.Atoi(params["command"])
	switch cmd {
	case 0:
		_, err = fmt.Fprint(pipe, "p")
	case 1:
		_, err = fmt.Fprint(pipe, "\x5b\x44")
	case 2:
		_, err = fmt.Fprint(pipe, "\x5b\x43")
	case 3:
		err = stop()
	case 4:
		_, err = fmt.Fprint(pipe, "\x5b\x42")
	case 5:
		_, err = fmt.Fprint(pipe, "\x5b\x41")
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, map[string]string{"error": err.Error()})
	}
}

func stop() error {
	var err error
	_, err = fmt.Fprint(pipe, "q")
	pipe = nil
	return err
}
