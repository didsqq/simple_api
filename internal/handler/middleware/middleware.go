package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", req.Method, req.URL.Path)

		next.ServeHTTP(w, req)

		log.Printf("Completed %s %s in %v", req.Method, req.URL.Path, time.Since(start))
	})
}
