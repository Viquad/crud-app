package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type authResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjEzMzY0NTUsImlhdCI6MTY2MTMzNTU1NSwianRpIjoiMyJ9.5LfGkxciCiJgEFV8yjX9Pvelt6sZtvUefgIiHIUiiak"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status" example:"ok"`
}

func newErrorResponse(c *gin.Context, statusCode int, context, problem string, err error) {
	logrus.WithFields(logrus.Fields{
		"context": context,
		"problem": problem,
	}).Error(err)
	c.AbortWithStatusJSON(statusCode, errorResponse{err.Error()})
}
