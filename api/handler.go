package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/simonjm/rasp-tv/data"
)

// Context struct that holds other objects that are needed by every handler
type Context struct {
	Logger *log.Logger
	Config *Config
	Db     data.RaspTvDataFetcher
}

// StreamHandlerFunc a custom handler function that includes a Context but doesn't return anything
type StreamHandlerFunc func(*Context, http.ResponseWriter, *http.Request)

// StreamHandler a Handler for the streaming routes
type StreamHandler struct {
	logger  *log.Logger
	config  *Config
	handler StreamHandlerFunc
}

// ServeHTTP implements the Handler interface for StreamHandler
func (a *StreamHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// initialize a connection to the database
	db, err := data.NewRaspTvDatabase(a.config.DbPath)
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

// HandlerFunc a custom handler function that includes a Context and returns a status code, the response body and an error
type HandlerFunc func(*Context, http.ResponseWriter, *http.Request) (int, interface{}, error)

// Handler a Handler for all routes that return data
type Handler struct {
	logger  *log.Logger
	config  *Config
	handler HandlerFunc
}

// ServeHTTP implements the Handler interface for ApiHandler
func (a *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// these routes always return JSON
	rw.Header().Add("Content-Type", "application/json")

	// initialize a connection to the database
	db, err := data.NewRaspTvDatabase(a.config.DbPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
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

	// calls the custom handler func
	status, data, err := a.handler(&context, rw, req)
	rw.WriteHeader(status)

	if err != nil {
		a.logger.Printf("[Error] - HTTP %d - %s\n", status, err)
		json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	// serialize and send back the result of the handler func
	json.NewEncoder(rw).Encode(data)
}

// NewAPIHandler constructs an ApiHandler
func NewAPIHandler(logger *log.Logger, config *Config, handler HandlerFunc) *Handler {
	return &Handler{
		logger:  logger,
		config:  config,
		handler: handler,
	}
}

// NewStreamHandler constructs a StreamHandler
func NewStreamHandler(logger *log.Logger, config *Config, handler StreamHandlerFunc) *StreamHandler {
	return &StreamHandler{
		logger:  logger,
		config:  config,
		handler: handler,
	}
}
