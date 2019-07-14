package middlewares

import (
	"fmt"
	"net/http"
)

// MiddleWare is a type thats used for custom middleware
type MiddleWare func(http.Handler) http.Handler

// captureStats logs each incoming request
func captureStats(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO fix this to capture formatted stats
		fmt.Println("Received a request")

		next.ServeHTTP(w, r)
	})
}

// TODO add more middle layers
