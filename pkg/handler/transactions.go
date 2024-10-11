package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) deposit(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	h.service.Transactions.Deposit(id)
}

func (h *Handler) withdraw(ctx *gin.Context) {

}

func (h *Handler) transfer(ctx *gin.Context) {

}

func (h *Handler) transactions(ctx *gin.Context) {

}
