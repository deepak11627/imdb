package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deepak11627/imdb/logger"

	"github.com/deepak11627/imdb/resources"
)

// MovieHandler handles Movie resource requests
type MovieHandler struct {
	service resources.MovieService
	logger  logger.Logger
}

func NewMovieHandler(s resources.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

func (m *MovieHandler) Handle(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid movie id %q", head), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		m.handleGet(id, res, req)
	case "PUT":
		m.handlePut(id, res, req)
	case "POST":
		m.handlePost(res, req)
	default:
		http.Error(res, "Only GET, POST and PUT are allowed", http.StatusMethodNotAllowed)
	}
}

func (mh *MovieHandler) handleGet(id int, res http.ResponseWriter, req *http.Request) {

}
func (mh *MovieHandler) handlePut(id int, res http.ResponseWriter, req *http.Request) {
}

type newMovieResponse struct {
	ID int `json:"id"`
}

func (mh *MovieHandler) handlePost(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var m resources.Movie
	err := decoder.Decode(&m)
	if err != nil {
		mh.logger.Warn("failed to create movie", "error", err)
		r, _ := jsonResponse(errorResponse{Code: strconv.Itoa(http.StatusBadRequest), Errors: []string{err.Error()}})
		res.Write(r)
		return
	}
	id, err := mh.service.CreateMovie(m.Name, m.Director, m.Genre, m.Score)
	if err != nil {
		mh.logger.Warn("failed to create movie", "error", err)
		r, _ := jsonResponse(errorResponse{Code: strconv.Itoa(http.StatusInternalServerError), Errors: []string{"internal server error"}})
		res.Write(r)
		return
	}
	r, err := jsonResponse(newMovieResponse{ID: id})
	if err != nil {
		mh.logger.Warn("failed to create movie json response", "error", err)
		return
	}
	res.Write(r)
}
