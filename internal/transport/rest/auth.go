package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) initAuth(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		// Register user
		auth.POST("/sign-up", h.signUp)
		// Authorization with access token
		auth.POST("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
		// Authorization with session
		auth.POST("/log-in", h.login)
		auth.GET("/log-out", h.logout)
	}
}

// SignUp godoc
// @Summary     SingUp
// @Description SingUp
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body     domain.SignUpInput true "user credentials to Sign-Up"
// @Success     201   {object} rest.statusResponse
// @Failure     400   {object} rest.errorResponse
// @Failure     409   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "SignUp()", "binding error", err)
		return
	}

	err := h.services.GetUserService().Create(c.Request.Context(), input)
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		newErrorResponse(c, http.StatusConflict, "SignUp()", "service error", err)
		return
	case err != nil:
		newErrorResponse(c, http.StatusInternalServerError, "SignUp()", "service error", err)
		return
	}

	c.JSON(http.StatusCreated, statusResponse{"ok"})
}

// SignIn godoc
// @Summary     SignIn
// @Description Get access token by user credentials
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body     domain.SignInInput true "user credentials to Sign-In"
// @Success     200   {object} rest.authResponse
// @Failure     400   {object} rest.errorResponse
// @Failure     404   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "SignIn()", "binding error", err)
		return
	}

	accessToken, refreshToken, err := h.services.GetUserService().GetTokenByCredentials(c.Request.Context(), input)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		newErrorResponse(c, http.StatusNotFound, "SignIn()", "user not found error", err)
		return
	case err != nil:
		newErrorResponse(c, http.StatusInternalServerError, "SignIn()", "service error", err)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.JSON(http.StatusCreated, authResponse{Token: accessToken})
}

// @Summary     Refresh tokens
// @Description Refresh tokens
// @Tags        auth
// @Produce     json
// @Success     201 {object} rest.authResponse
// @Failure     400   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /auth/refresh [get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh-token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "refresh()", "cookie error", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"context": "refresh()",
	}).Debugf("%s", cookie.Value)

	accessToken, refreshToken, err := h.services.GetUserService().RefreshTokens(c.Request.Context(), cookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "refresh()", "service error", err)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.JSON(http.StatusCreated, authResponse{Token: accessToken})
}

// @Summary     LogIn
// @Description Init session by user credentials
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body     domain.SignInInput true "user credentials to Log-In"
// @Success     200   {object} rest.statusResponse
// @Failure     400 {object} rest.errorResponse
// @Failure     404   {object} rest.errorResponse
// @Failure     500 {object} rest.errorResponse
// @Router      /auth/log-in [post]
func (h *Handler) login(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "login()", "binding error", err)
		return
	}

	err := h.services.GetUserService().InitSession(c, input)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		newErrorResponse(c, http.StatusNotFound, "login()", "user not found error", err)
		return
	case err != nil:
		newErrorResponse(c, http.StatusInternalServerError, "login()", "service error", err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary     LogOut
// @Description Drop user's session
// @Tags        auth
// @Produce     json
// @Success     200 {object} rest.statusResponse
// @Failure     500 {object} rest.errorResponse
// @Router      /auth/log-out [get]
func (h *Handler) logout(c *gin.Context) {
	err := h.services.GetUserService().DropSession(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "logout()", "service error", err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
