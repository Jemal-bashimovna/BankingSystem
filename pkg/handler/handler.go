package handler

import (
	"bankingsystem/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api/accounts")
	{
		api.POST("", h.createAccount)               //
		api.DELETE("/:account_id", h.deleteAccount) //
		api.GET("/:account_id", h.getAccount)       //
		api.GET("", h.getAllAccounts)               //

		api.POST("/:account_id/deposit", h.deposit)
		api.POST("/:account_id/withdraw", h.withdraw)
		api.POST("/:account_id/transfer", h.transfer)
		api.GET("/:account_id/transactions", h.transactions)
	}
	return router
}
