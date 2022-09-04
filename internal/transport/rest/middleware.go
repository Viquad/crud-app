package rest

import (
	"context"
	"errors"
	"fmt"
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
	err1 := h.authWithToken(c)
	if err1 == nil {
		c.Next()
		return
	}

	err2 := h.authWithSession(c)
	if err2 == nil {
		c.Next()
		return
	}

	err := fmt.Errorf("token: %s; session: %s", err1.Error(), err2.Error())
	newErrorResponse(c, http.StatusUnauthorized, "authMiddleware", "authorize error", err)
}

func (h *Handler) authWithToken(c *gin.Context) error {
	token, err := getTokenFromRequest(c)
	if err != nil {
		return err
	}

	userId, err := h.services.GetUserService().ParseToken(c.Request.Context(), token)
	if err != nil {
		return err
	}

	withUserId(c, userId)

	return nil
}

func (h *Handler) authWithSession(c *gin.Context) error {
	userId, err := h.services.GetUserService().GetSession(c)
	if err != nil {
		return err
	}

	withUserId(c, userId)

	return nil
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

func withUserId(c *gin.Context, userId int64) {
	ctx := context.WithValue(c.Request.Context(), domain.UserIdKey, userId)
	c.Request = c.Request.WithContext(ctx)
}
