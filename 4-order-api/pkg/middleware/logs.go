package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		latency := time.Since(start)
		entry := logger.WithFields(logrus.Fields{
			"status":     wrapper.StatusCode,
			"method":     r.Method,
			"path":       r.URL.Path,
			"query":      r.URL.RawQuery,
			"remote_ip":  r.RemoteAddr,
			"user_agent": r.UserAgent(),
			"latency_ms": latency.Milliseconds(),
		})

		if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
			entry = entry.WithField("request_id", requestID)
		}

		switch {
		case wrapper.StatusCode >= 500:
			entry.Error("http request completed")
		case wrapper.StatusCode >= 400:
			entry.Warn("http request completed")
		default:
			entry.Info("http request completed")
		}
	})
}
