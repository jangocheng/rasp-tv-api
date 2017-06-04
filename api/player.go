package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/simonjm/rasp-tv-api/data"

	"github.com/gorilla/mux"
)

// the possible commands to send to omxplayer
const (
	TOGGLE = iota
	BACKWARD
	FORWARD
	STOP
	FASTBACKWARD
	FASTFORWARD
)

var pipe io.WriteCloser

// startPlayer calls omxplayer with the given path and returns the PID
func startPlayer(path string) (int64, error) {
	// stop any currently playing videos
	var err error
	if err = stop(); err != nil {
		return -1, err
	}

	command := exec.Command("omxplayer.bin", "-o", "hdmi", "-b", path)
	// set up pip to stdin so we can send commands to control the playback
	pipe, err = command.StdinPipe()
	if err != nil {
		return -1, err
	}

	// starts the player
	err = command.Start()
	go func() {
		// waits for the command in a goroutine so the process gets cleaned up correctly
		command.Wait()
		pipe = nil
	}()

	return int64(command.Process.Pid), err
}

// RunPlayerCommand route that sends a command the omxplayer process to control the playback
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

	db := context.Db
	switch cmd {
	case TOGGLE:
		session, e := db.GetSession()
		if e != nil {
			err = e
			break
		}
		// toggles the IsPaused value from the session and save the new value
		session.IsPaused = !session.IsPaused
		if e = db.SaveSession(session); e != nil {
			err = e
			break
		}
		// pause the player
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

// NowPlaying route that returns the current session
func NowPlaying(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	session, err := context.Db.GetSession()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if session.ID == 0 {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, session, nil
}

// ClearSession route that clears the current session
func ClearSession(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	if err := context.Db.ClearSession(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, statusResponse("Success"), nil
}

// UpdateSession route that updates the session. The data to save is deserialized from the request body
func UpdateSession(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	defer req.Body.Close()
	session := &data.Session{}

	if err := json.NewDecoder(req.Body).Decode(session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	db := context.Db
	if err := db.ClearSession(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := db.SaveSession(session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, session, nil
}

// stop sends the quit command the omxplayer process and clears the pipe
func stop() error {
	var err error
	if pipe != nil {
		_, err = fmt.Fprint(pipe, "q")
		pipe = nil
	}

	return err
}
