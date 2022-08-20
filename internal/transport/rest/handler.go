package rest

import (
	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
)

type Services interface {
	GetAccountService() domain.AccountService
}

type Handler struct {
	services Services
}

func NewHandler(s Services) *Handler {
	return &Handler{s}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(Logger(), gin.Recovery())

	h.initAccount(&router.RouterGroup)

	return router
}
