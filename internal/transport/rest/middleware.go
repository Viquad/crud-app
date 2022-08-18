package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"url":    c.Request.URL.String(),
		}).Info("Accept request")

		t := time.Now()

		c.Next()

		logrus.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"elapsed": time.Since(t).String(),
		}).Info("Send response")
	}
}
