package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bombsimon/laundry/errors"
	"github.com/bombsimon/laundry/log"
	"github.com/sirupsen/logrus"
)

// loggingResponseWriter embeds http.ResponseWriter and add a field for status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	err        *errors.LaundryError
}

// NewLoggingResponseWriter will return a new respone writer capable of tacking status code
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, nil}
}

// WriteHeader overrides the default http.WriteHeader to store the written status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// WriteError ...
func (lrw *loggingResponseWriter) WriteError(err error) {
	lrw.err = err.(*errors.LaundryError)
}

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := NewLoggingResponseWriter(w)
			h.ServeHTTP(lrw, r)

			defer log.GetLogger().WithFields(logrus.Fields{
				"method":         r.Method,
				"remote address": r.RemoteAddr,
				"path":           r.URL.String(),
				"protocol":       r.Proto,
				"content length": r.ContentLength,
				"status":         http.StatusText(lrw.statusCode),
				"status code":    lrw.statusCode,
				"elapsed":        fmt.Sprintf("%02f %s", time.Now().Sub(startTime).Seconds()*1000, "ms"),
			}).Infof("Request handled")
		})
	}
}
