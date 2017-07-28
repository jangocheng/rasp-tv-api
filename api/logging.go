package api

import (
	"encoding/json"
	"net/http"

	"github.com/simonjm/rasp-tv-api/data"
)

// SaveLogs route for saving a batch of log entries. The data to save is deserialized from the request body
func SaveLogs(context *Context, rw http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	var logs []data.Log
	if err := json.NewDecoder(req.Body).Decode(&logs); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if len(logs) == 0 {
		return http.StatusOK, statusResponse("No logs to save"), nil
	}

	if err := context.Db.SaveLogs(logs); err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, statusResponse("Logs were successfully inserted"), nil
}
