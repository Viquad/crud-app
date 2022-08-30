package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Logger(c *gin.Context) {
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

func (h *Handler) authMiddleware(c *gin.Context) {
	token, err := getTokenFromRequest(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "authMiddleware", "get token error", err)
		return
	}

	userId, err := h.services.GetUserService().ParseToken(c.Request.Context(), token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "authMiddleware", "service error", err)
	}

	ctx := context.WithValue(c.Request.Context(), domain.UserIdKey, userId)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
