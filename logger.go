package kitchen

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

// LoggerMiddleware is a simple middleware to handle logging with logrus
func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logrus.WithFields(logrus.Fields{
			"method":  r.Method,
			"request": r.RequestURI,
			"remote":  r.RemoteAddr,
		}).Info("started handling request")

		go func() {
			res := w.(ResponseWriter)
			<-res.Context().Done()
			latency := time.Since(start)
			logrus.WithFields(logrus.Fields{
				"status":          res.Status(),
				"method":          r.Method,
				"request":         r.RequestURI,
				"remote":          r.RemoteAddr,
				"text_status":     http.StatusText(res.Status()),
				"took":            latency,
				"measure_latency": latency.Nanoseconds(),
				"cancel_error":    res.Context().Err(),
			}).Info("completed handling request")
		}()

		next.ServeHTTP(w, r) // Call the next one
	}
	return http.HandlerFunc(fn)
}
