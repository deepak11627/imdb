package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/deepak11627/imdb/db"
	logger "github.com/deepak11627/imdb/logger"
)

// Handler handles all incoming requests
type Handler struct {
	Database db.DatabaseService
	logger   logger.Logger
	// We could use http.Handler as a type here; using the specific type has
	// the advantage that static analysis tools can link directly from
	// h.MovieHandler.ServeHTTP to the correct definition. The disadvantage is
	// that we have slightly stronger coupling.
	MovieHandler *MovieHandler
}

// NewHandler returns the Custom Hanlder type
func NewHandler(db db.DatabaseService, l logger.Logger) *Handler {
	return &Handler{
		Database:     db,
		logger:       l,
		MovieHandler: &MovieHandler{service: db, logger: l},
	}
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "movie":
		h.MovieHandler.Handle(res, req)
		return
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}

}

type errorResponse struct {
	code   int
	Errors []string `json:"errors"`
}

func jsonResponse(res http.ResponseWriter, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return err
	}
	res.Header().Set("Content-Type", "application/json")
	if e, isError := v.(errorResponse); isError {
		res.WriteHeader(e.code)
	} else {
		res.WriteHeader(http.StatusOK)
	}
	res.Write(b)
	return nil
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
