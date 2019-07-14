package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/deepak11627/imdb/db"
)

// Handler handles all incoming requests
type Handler struct {
	Database db.DatabaseService

	// We could use http.Handler as a type here; using the specific type has
	// the advantage that static analysis tools can link directly from
	// h.MovieHandler.ServeHTTP to the correct definition. The disadvantage is
	// that we have slightly stronger coupling.
	MovieHandler *MovieHandler
	//MiddleWares  []middlewares.MiddleWare
}

// NewHandler returns the Custom Hanlder type
func NewHandler(db db.DatabaseService) *Handler {
	return &Handler{
		Database:     db,
		MovieHandler: &MovieHandler{service: db},
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
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}

func jsonResponse(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

// type Handler interface {
// 	Get(int)
// 	Post()
// 	Put(int)
// 	Delete(int)
// }

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
