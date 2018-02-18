package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bombsimon/laundry/errors"
	"github.com/bombsimon/laundry/log"
	"github.com/sirupsen/logrus"
)

// LoggingResponseWriter embeds http.ResponseWriter and add a field for status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	err        *errors.LaundryError
}

// NewLoggingResponseWriter will return a new respone writer capable of tacking status code
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, nil}
}

// WriteHeader overrides the default http.WriteHeader to store the written status code
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// WriteError can be used to store an error from the handler in the response writer
func (lrw *LoggingResponseWriter) WriteError(err error) {
	switch e := err.(type) {
	case *errors.LaundryError:
		lrw.err = e
	case error:
		lrw.err = errors.New(e)
	}
}

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := NewLoggingResponseWriter(w)
			h.ServeHTTP(lrw, r)

			defer func() {
				logger := log.GetLogger().WithFields(logrus.Fields{
					"method":         r.Method,
					"remote address": r.RemoteAddr,
					"path":           r.URL.String(),
					"protocol":       r.Proto,
					"content length": r.ContentLength,
					"status":         http.StatusText(lrw.statusCode),
					"status code":    lrw.statusCode,
					"elapsed":        fmt.Sprintf("%.3f %s", time.Now().Sub(startTime).Seconds()*1000, "ms"),
				})

				if lrw.err != nil {
					logger.WithFields(logrus.Fields{
						"error": lrw.err.Reasons,
					}).Errorf("Error processing request")
					return
				}

				logger.Infof("Request processed")
			}()
		})
	}
}
