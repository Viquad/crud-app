package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAccount(router *gin.RouterGroup) {
	account := router.Group("/account")
	{
		account.Use(h.authMiddleware)

		account.POST("/", h.CreateAccount)
		account.PUT("/", h.CreateAccount)
		account.GET("/", h.GetAccounts)
		account.GET("/:id", h.GetAccountById)
		account.POST("/:id", h.UpdateAccount)
		account.PUT("/:id", h.UpdateAccount)
		account.DELETE("/:id", h.DeleteAccount)
	}
}

// CreateAccount godoc
// @Summary     Create new account for user
// @Description Create new account for user
// @Security    ApiKeyAuth
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       input body     domain.AccountCreateInput true "account info"
// @Success     200   {object} domain.Account
// @Failure     400   {object} rest.errorResponse
// @Failure     401   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /account [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var input domain.AccountCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "CreateAccount()", "binding error", err)
		return
	}

	account, err := h.services.GetAccountService().Create(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "CreateAccount()", "service error", err)
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccountById godoc
// @Summary     Get account
// @Description Get user's account by id
// @Security    ApiKeyAuth
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id  path     string true "account id"
// @Success     200 {object} domain.Account
// @Failure     400 {object} rest.errorResponse
// @Failure     401 {object} rest.errorResponse
// @Failure     404 {object} rest.errorResponse
// @Failure     500 {object} rest.errorResponse
// @Router      /account/{id} [get]
func (h *Handler) GetAccountById(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "GetAccountById()", "parsing id error", err)
		return
	}

	account, err := h.services.GetAccountService().GetById(c.Request.Context(), id)
	if err != nil {
		context, problem := "GetAccountById()", "service error"
		switch {
		case errors.Is(err, domain.ErrNotExist):
			newErrorResponse(c, http.StatusNotFound, context, problem, err)
		default:
			newErrorResponse(c, http.StatusInternalServerError, context, problem, err)
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccounts godoc
// @Summary     Get accounts
// @Description Get all user's accounts list
// @Security    ApiKeyAuth
// @Tags        account
// @Produce     json
// @Success     200 {object} []domain.Account
// @Failure     401 {object} rest.errorResponse
// @Failure     500 {object} rest.errorResponse
// @Router      /account [get]
func (h *Handler) GetAccounts(c *gin.Context) {
	accounts, err := h.services.GetAccountService().List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "GetAccounts()", "service error", err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount godoc
// @Summary     Update account
// @Description Update user's account info by id
// @Security    ApiKeyAuth
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id    path     string                    true "account id"
// @Param       input body     domain.AccountUpdateInput true "account update info"
// @Success     200   {object} domain.Account
// @Failure     400   {object} rest.errorResponse
// @Failure     401   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /account/{id} [post]
func (h *Handler) UpdateAccount(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "parsing id error", err)
		return
	}

	var input domain.AccountUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "binding error", err)
		return
	}

	account, err := h.services.GetAccountService().UpdateById(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUpdateFailed):
			newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "service error", err)
		default:
			newErrorResponse(c, http.StatusInternalServerError, "UpdateAccount()", "service error", err)
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary     Delete account
// @Description Delete user's account by id
// @Security    ApiKeyAuth
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id  path     string true "account id"
// @Success     200 {object} rest.statusResponse
// @Failure     400 {object} rest.errorResponse
// @Failure     401 {object} rest.errorResponse
// @Failure     404 {object} rest.errorResponse
// @Failure     500 {object} rest.errorResponse
// @Router      /account/{id} [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "DeleteAccount()", "parsing id error", err)
		return
	}

	if err := h.services.GetAccountService().DeleteById(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domain.ErrDeleteFailed):
			newErrorResponse(c, http.StatusNotFound, "DeleteAccount()", "service error", err)
		default:
			newErrorResponse(c, http.StatusInternalServerError, "DeleteAccount()", "service error", err)
		}
		return
	}

	c.JSON(http.StatusOK, statusResponse{"OK"})
}

func parseId(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return id, err
	}

	if id < 1 {
		return id, domain.ErrInvalidId
	}

	return id, err
}
