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
// @Summary     Create new account
// @Description Create new account
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       input body     domain.Account true "account info"
// @Success     200   {object} domain.Account
// @Failure     400   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /account [put]
// @Router      /account [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	account := new(domain.Account)

	if err := c.ShouldBindJSON(account); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "CreateAccount()", "binding error", err.Error())
		return
	}

	account, err := h.services.GetAccountService().Create(c.Request.Context(), *account)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "CreateAccount()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccountById godoc
// @Summary     Get account
// @Description Get account by id
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id      path     string true "account id"
// @Success     200     {object} domain.Account
// @Failure     400,404 {object} rest.errorResponse
// @Failure     500     {object} rest.errorResponse
// @Router      /account/{id} [get]
func (h *Handler) GetAccountById(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "GetAccountById()", "parsing id error", err.Error())
		return
	}

	account, err := h.services.GetAccountService().GetById(c.Request.Context(), id)
	if err != nil {
		context, problem := "GetAccountById()", "service error"
		switch {
		case errors.Is(err, domain.ErrNotExist):
			newErrorResponse(c, http.StatusNotFound, context, problem, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, context, problem, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccounts godoc
// @Summary     Get accounts
// @Description Get all accounts list
// @Tags        account
// @Accept      json
// @Produce     json
// @Success     200 {object} []domain.Account
// @Failure     500 {object} rest.errorResponse
// @Router      /account [get]
func (h *Handler) GetAccounts(c *gin.Context) {
	accounts, err := h.services.GetAccountService().List(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "GetAccounts()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount godoc
// @Summary     Update account
// @Description Update account info by id
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id    path     string                    true "account id"
// @Param       input body     domain.AccountUpdateInput true "account update info"
// @Success     200   {object} domain.Account
// @Failure     400   {object} rest.errorResponse
// @Failure     500   {object} rest.errorResponse
// @Router      /account/{id} [put]
// @Router      /account/{id} [post]
func (h *Handler) UpdateAccount(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "parsing id error", err.Error())
		return
	}

	var input domain.AccountUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "binding error", err.Error())
		return
	}

	account, err := h.services.GetAccountService().UpdateById(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUpdateFailed):
			newErrorResponse(c, http.StatusBadRequest, "UpdateAccount()", "service error", err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, "UpdateAccount()", "service error", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary     Delete account
// @Description Delete account by id
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id      path     string true "account id"
// @Success     200     {object} rest.statusResponse
// @Failure     400,404 {object} rest.errorResponse
// @Failure     500     {object} rest.errorResponse
// @Router      /account/{id} [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "DeleteAccount()", "parsing id error", err.Error())
		return
	}

	if err := h.services.GetAccountService().DeleteById(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domain.ErrDeleteFailed):
			newErrorResponse(c, http.StatusNotFound, "DeleteAccount()", "service error", err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, "DeleteAccount()", "service error", err.Error())
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
