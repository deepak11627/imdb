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

	switch req.Method {
	case "GET":
		m.handleGet(res, req)
	case "PUT":
		m.handlePut(res, req)
	case "POST":
		m.handlePost(res, req)
	default:
		http.Error(res, "Only GET, POST and PUT are allowed", http.StatusMethodNotAllowed)
	}
}

type movieResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name,omitempty"`
	Score    float32 `json:"imdb_score,omitempty"`
	Director string  `json:"director,omitempty"`
}

func (mh *MovieHandler) handleGet(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid movie id %q", head), http.StatusBadRequest)
		return
	}
	m, err := mh.service.GetMovie(id)
	if err != nil {
		r, _ := jsonResponse(errorResponse{Code: strconv.Itoa(http.StatusInternalServerError), Errors: []string{err.Error()}})
		res.Write(r)
		return
	}
	r, err := jsonResponse(movieResponse{ID: m.ID, Name: m.Name, Score: m.Score, Director: m.Director})
	if err != nil {
		mh.logger.Warn("failed to create movie json response", "error", err)
		return
	}
	res.Write(r)
}
func (mh *MovieHandler) handlePut(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	_, err := strconv.Atoi(head)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid movie id %q", head), http.StatusBadRequest)
		return
	}

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
	r, err := jsonResponse(movieResponse{ID: id})
	if err != nil {
		mh.logger.Warn("failed to create movie json response", "error", err)
		return
	}
	res.Write(r)
}
