package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func serverLogger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			log.WithFields(log.Fields{
				"method": method,
				"path":   path,
				"code":   res.Status(),
				"time":   stop.Sub(start).String(),
			}).Debug("server request")
			return nil
		}
	}
}
