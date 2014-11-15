package kitchen

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logrus.WithFields(logrus.Fields{
			"method":  r.Method,
			"request": r.RequestURI,
			"remote":  r.RemoteAddr,
		}).Info("started handling request")

		next.ServeHTTP(w, r) // Call the next one

		latency := time.Since(start)
		if res, ok := w.(ResponseWriter); ok {
			logrus.WithFields(logrus.Fields{
				"status":          res.Status(),
				"method":          r.Method,
				"request":         r.RequestURI,
				"remote":          r.RemoteAddr,
				"text_status":     http.StatusText(res.Status()),
				"took":            latency,
				"measure latency": latency.Nanoseconds(),
			}).Info("completed handling request")
		} else {
			logrus.WithFields(logrus.Fields{
				"status":           "NOT SUPPORT THE kitchen ",
				"method":           r.Method,
				"request":          r.RequestURI,
				"remote":           r.RemoteAddr,
				"text_status":      "NOT SUPPORT THE kitchen ",
				"took":             latency,
				"measure #latency": latency.Nanoseconds(),
			}).Info("completed handling request")
		}
	}
	return http.HandlerFunc(fn)
}
