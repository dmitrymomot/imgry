package server

import (
	"fmt"
	"net/http"

	"github.com/tobi/airbrake-go"
)

// Airbrake recoverer middleware to capture and report any panics to
// airbrake.io.
func AirbrakeRecoverer(apiKey string) func(http.Handler) http.Handler {
	airbrake.ApiKey = apiKey
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey != "" {
				defer airbrake.CapturePanic(r)
			}
			next.ServeHTTP(w, r)
		})
	}
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *wrappedResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *wrappedResponseWriter) Status() int {
	return l.status
}

func trackRoute(metricID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		route := fmt.Sprintf("route.%s", metricID)
		errRoute := fmt.Sprintf("%s-err", route)

		handler := func(w http.ResponseWriter, r *http.Request) {
			lw := &wrappedResponseWriter{w, -1}
			next.ServeHTTP(lw, r)
		}
		return http.HandlerFunc(handler)
	}
}
