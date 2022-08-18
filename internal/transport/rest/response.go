package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, context, problem, message string) {
	logrus.WithFields(logrus.Fields{
		"context": context,
		"problem": problem,
	}).Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
