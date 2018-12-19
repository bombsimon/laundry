package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingResponseWriter embeds http.ResponseWriter and add a field for status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	err        error
}

// NewLoggingResponseWriter will return a new response writer capable of
// tacking status code
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
	lrw.err = err
}

// Logger is a semi verbose JSON logger that will log all requests.
func Logger() Adapter {
	logger := logrus.New()

	return func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := NewLoggingResponseWriter(w)
			h.ServeHTTP(lrw, r)

			defer func() {
				fieldLogger := logger.WithFields(logrus.Fields{
					"method":         r.Method,
					"remote_address": r.RemoteAddr,
					"path":           r.URL.String(),
					"protocol":       r.Proto,
					"content_length": r.ContentLength,
					"status":         http.StatusText(lrw.statusCode),
					"status_code":    lrw.statusCode,
					"elapsed":        fmt.Sprintf("%.3f %s", time.Now().Sub(startTime).Seconds()*1000, "ms"),
				})

				if lrw.err != nil {
					fieldLogger.WithFields(logrus.Fields{
						"error": lrw.err.Error(),
					}).Errorf("Error processing request")
					return
				}

				fieldLogger.Infof("Request processed")
			}()
		})
	}
}
