package rest

import (
	"github.com/Viquad/crud-app/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	accountService domain.AccountService
}

func NewHandler(s domain.Services) *Handler {
	return &Handler{
		accountService: s.GetAccountService(),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(Logger(), gin.Recovery())

	h.initAccount(&router.RouterGroup)

	return router
}
