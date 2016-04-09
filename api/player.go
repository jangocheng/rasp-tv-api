package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"simongeeks.com/joe/rasp-tv/data"

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

func startPlayer(path string) (int64, error) {
	var err error
	if err = stop(); err != nil {
		return -1, err
	}

	command := exec.Command("omxplayer.bin", "-o", "hdmi", "-b", path)
	pipe, err = command.StdinPipe()
	if err != nil {
		return -1, err
	}

	err = command.Start()
	go func() {
		command.Wait()
		pipe = nil
	}()

	return int64(command.Process.Pid), err
}

func RunPlayerCommand(r render.Render, params martini.Params, db *sql.DB, logger *log.Logger) {
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
		return
	}

	switch cmd {
	case TOGGLE:
		session, e := data.GetSession(db)
		if e != nil {
			err = e
			break
		}
		session.IsPaused = !session.IsPaused
		if e = session.Save(db); e != nil {
			err = e
			break
		}
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

func NowPlaying(r render.Render, db *sql.DB, logger *log.Logger) {
	session, err := data.GetSession(db)
	if err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if session.Id == 0 {
		return
	}

	r.JSON(200, session)
}

func ClearSession(r render.Render, db *sql.DB, logger *log.Logger) {
	if err := data.ClearSessions(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, statusResponse("Success"))
}

func UpdateSession(r render.Render, req *http.Request, db *sql.DB, logger *log.Logger) {
	defer req.Body.Close()
	session := &data.Session{}

	if err := json.NewDecoder(req.Body).Decode(session); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if err := data.ClearSessions(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	if err := session.Save(db); err != nil {
		logger.Println(errorMsg(err.Error()))
		r.JSON(500, errorResponse(err))
		return
	}

	r.JSON(200, session)
}

func stop() error {
	var err error
	if pipe != nil {
		_, err = fmt.Fprint(pipe, "q")
		pipe = nil
	}

	return err
}
