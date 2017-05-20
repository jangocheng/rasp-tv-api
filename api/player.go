package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/simonjm/rasp-tv/data"

	"github.com/gorilla/mux"
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

func RunPlayerCommand(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var err error
	if pipe == nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("Player not started")
	}

	cmdStr := mux.Vars(req)["command"]
	cmd, err := strconv.Atoi(cmdStr)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("Invalid command: %s", cmdStr)
	}

	switch cmd {
	case TOGGLE:
		session, e := data.GetSession(context.Db)
		if e != nil {
			err = e
			break
		}
		session.IsPaused = !session.IsPaused
		if e = session.Save(context.Db); e != nil {
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
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Success"), nil
}

func NowPlaying(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	session, err := data.GetSession(context.Db)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if session.Id == 0 {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, session, nil
}

func ClearSession(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	if err := data.ClearSessions(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Success"), nil
}

func UpdateSession(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	defer req.Body.Close()
	session := &data.Session{}

	if err := json.NewDecoder(req.Body).Decode(session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := data.ClearSessions(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := session.Save(context.Db); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, session, nil
}

func stop() error {
	var err error
	if pipe != nil {
		_, err = fmt.Fprint(pipe, "q")
		pipe = nil
	}

	return err
}
