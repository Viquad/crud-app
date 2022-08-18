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
		account.POST("/", h.CreateAccount)
		account.PUT("/", h.CreateAccount)
		account.GET("/", h.GetAccounts)
		account.GET("/:id", h.GetAccountById)
		account.POST("/:id", h.UpdateAccount)
		account.PUT("/:id", h.UpdateAccount)
		account.DELETE("/:id", h.DeleteAccount)
	}
}

func (h *Handler) CreateAccount(c *gin.Context) {
	account := new(domain.Account)

	if err := c.ShouldBindJSON(account); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "CreateAccount()", "binding error", err.Error())
		return
	}

	account, err := h.accountService.Create(c.Request.Context(), *account)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "CreateAccount()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccountById(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "GetAccountById()", "parsing id error", err.Error())
		return
	}

	account, err := h.accountService.GetById(c.Request.Context(), id)
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

func (h *Handler) GetAccounts(c *gin.Context) {
	accounts, err := h.accountService.All(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "GetAccounts()", "service error", err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

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

	account, err := h.accountService.Update(c.Request.Context(), id, input)
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

func (h *Handler) DeleteAccount(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "DeleteAccount()", "parsing id error", err.Error())
		return
	}

	if err := h.accountService.Delete(c.Request.Context(), id); err != nil {
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
 