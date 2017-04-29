package middlewares

import (
	"net/http"
	"time"

	echo "gopkg.in/labstack/echo.v3"

	"github.com/Sirupsen/logrus"
)

// Logger is the middleware for log system
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start timer
		start := time.Now()
		// Process request
		err := next(c)

		latency := time.Since(start)
		statusCode := c.Response().Status
		logrus.WithFields(
			logrus.Fields{
				"Method":   c.Request().Method,
				"Path":     c.Request().URL.Path,
				"Latency":  latency,
				"ClientIP": c.RealIP(),
				"Status":   statusCode,
				"Err":      err,
			},
		).Debug(http.StatusText(statusCode))

		return err
	}
}
