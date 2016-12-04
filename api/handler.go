package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	Logger *log.Logger
	Config *Config
	Db     *sql.DB
}

type StreamHandlerFunc func(*Context, http.ResponseWriter, *http.Request)

type StreamHandler struct {
	logger  *log.Logger
	config  *Config
	handler StreamHandlerFunc
}

func (a *StreamHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("sqlite3", a.config.DbPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		a.logger.Printf("[Error] - %s\n", err)
		return
	}
	defer db.Close()

	context := Context{
		Logger: a.logger,
		Config: a.config,
		Db:     db,
	}

	a.handler(&context, rw, req)
}

type ApiHandlerFunc func(*Context, http.ResponseWriter, *http.Request) (int, interface{}, error)

type ApiHandler struct {
	logger  *log.Logger
	config  *Config
	handler ApiHandlerFunc
}

func (a *ApiHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("sqlite3", a.config.DbPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Add("Content-Type", "application/json")
		a.logger.Printf("[Error] - %s\n", err)
		json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer db.Close()

	context := Context{
		Logger: a.logger,
		Config: a.config,
		Db:     db,
	}

	status, data, err := a.handler(&context, rw, req)
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")

	if err != nil {
		a.logger.Printf("[Error] - HTTP %d - %s\n", status, err)
		json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(rw).Encode(data)
}

func NewApiHandler(logger *log.Logger, config *Config, handler ApiHandlerFunc) *ApiHandler {
	return &ApiHandler{
		logger:  logger,
		config:  config,
		handler: handler,
	}
}

func NewStreamHandler(logger *log.Logger, config *Config, handler StreamHandlerFunc) *StreamHandler {
	return &StreamHandler{
		logger:  logger,
		config:  config,
		handler: handler,
	}
}
