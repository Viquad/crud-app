package rest

import (
	"errors"
	"net/http"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAuth(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
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
// @Failure     500   {object} rest.errorResponse
// @Router      /auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "SignUp()", "binding error", err.Error())
		return
	}

	err := h.services.GetUserService().Create(c.Request.Context(), input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "SignUp()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, statusResponse{"ok"})
}

// SignIn godoc
// @Summary     SignIn
// @Description SignIn
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body     domain.SignInInput true "user credentials to Sign-In"
// @Success     200   {object} rest.userResponse
// @Failure     400   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "SignIn()", "binding error", err.Error())
		return
	}

	token, err := h.services.GetUserService().GetTokenByCredentials(c.Request.Context(), input)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		newErrorResponse(c, http.StatusNotFound, "SignIn()", "user not found error", err.Error())
		return
	case err != nil:
		newErrorResponse(c, http.StatusInternalServerError, "SignIn()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}
