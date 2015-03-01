package api

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

const (
	TOGGLE = iota
	BACKWARD
	FORWARD
	STOP
	FASTBACKWARD
	FASTFORWARD
)

var pipe io.WriteCloser

func startPlayer(path string) error {
	var err error
	if err = stop(); err != nil {
		return err
	}

	command := exec.Command("omxplayer", "-o", "hdmi", "-b", path)
	pipe, err = command.StdinPipe()
	if err != nil {
		return err
	}

	err = command.Start()
	go func() {
		command.Wait()
		pipe = nil
	}()

	return err
}

func RunPlayerCommand(r render.Render, params martini.Params, logger *log.Logger) {
	var err error
	if pipe == nil {
		err = fmt.Errorf("Player not started")
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	cmd, err := strconv.Atoi(params["command"])
	if err != nil {
		err = fmt.Errorf("Invalid command: %s", params["command"])
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
	}

	switch cmd {
	case TOGGLE:
		_, err = fmt.Fprint(pipe, "p")
	case BACKWARD:
		_, err = fmt.Fprint(pipe, "\x5b\x44")
	case FORWARD:
		_, err = fmt.Fprint(pipe, "\x5b\x43")
	case STOP:
		err = stop()
	case FASTBACKWARD:
		_, err = fmt.Fprint(pipe, "\x5b\x42")
	case FASTFORWARD:
		_, err = fmt.Fprint(pipe, "\x5b\x41")
	}

	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
	}
}

func stop() error {
	var err error
	if pipe != nil {
		_, err = fmt.Fprint(pipe, "q")
		pipe = nil
	}

	return err
}
