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

// TODO add Option type for logger

// NewMovieHandler create a movie handler
func NewMovieHandler(s resources.MovieService, l logger.Logger) *MovieHandler {
	return &MovieHandler{service: s, logger: l}
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
	ID       int      `json:"id"`
	Name     string   `json:"name,omitempty"`
	Score    string   `json:"imdb_score,omitempty"`
	Director string   `json:"director,omitempty"`
	Genre    []string `json:"genre,omitempty"`
}

func (mh *MovieHandler) handleGet(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{fmt.Sprintf("Invalid movie id %q", head)}})
		return
	}
	m, err := mh.service.GetMovie(id)
	if err != nil {
		mh.logger.Warn("error while getting movie", "error", err)
		jsonResponse(res, errorResponse{code: http.StatusInternalServerError, Errors: []string{err.Error()}})
		return
	}
	if m.IsEmpty() {
		jsonResponse(res, errorResponse{code: http.StatusNotFound, Errors: []string{"not found"}})
		return
	}
	err = jsonResponse(res, movieResponse{ID: m.ID, Name: m.Name, Score: m.Score, Director: m.Director, Genre: m.Genre})
	if err != nil {
		mh.logger.Warn("failed to create movie json response", "error", err)
		return
	}

}
func (mh *MovieHandler) handlePut(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{fmt.Sprintf("Invalid movie id %q", head)}})
		return
	}

	decoder := json.NewDecoder(req.Body)
	var m resources.Movie
	err = decoder.Decode(&m)
	if err != nil {
		mh.logger.Warn("failed to update movie", "error", err)
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{err.Error()}})
		return
	}
	if !m.IsValid() {
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{"Bad Data"}})
		return
	}

	mh.logger.Debug("received update movie request", m)
	err = mh.service.SaveMovie(id, m.Name, m.Director, m.Genre, m.Score)
	if err != nil {
		mh.logger.Warn("failed to update movie", "error", err)
		jsonResponse(res, errorResponse{code: http.StatusInternalServerError, Errors: []string{"internal server error"}})
		return
	}
}

func (mh *MovieHandler) handlePost(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var m resources.Movie
	err := decoder.Decode(&m)
	if err != nil {
		mh.logger.Warn("failed to create movie", "error", err)
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{err.Error()}})

		return
	}
	if !m.IsValid() {
		jsonResponse(res, errorResponse{code: http.StatusBadRequest, Errors: []string{"Bad Data"}})
		return
	}

	mh.logger.Debug("received create movie request", m)
	id, err := mh.service.CreateMovie(m.Name, m.Director, m.Genre, m.Score)
	if err != nil {
		mh.logger.Warn("failed to create movie", "error", err)
		jsonResponse(res, errorResponse{code: http.StatusInternalServerError, Errors: []string{"internal server error"}})
		return
	}
	err = jsonResponse(res, movieResponse{ID: id})
	if err != nil {
		mh.logger.Warn("failed to create movie json response", "error", err)
		return
	}

}
