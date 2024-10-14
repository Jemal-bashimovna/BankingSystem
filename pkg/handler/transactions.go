package handler

import (
	"bankingsystem/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) deposit(ctx *gin.Context) {

	var sum models.InputDeposit

	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&sum); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.DepositProducer(id, sum)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) withdraw(ctx *gin.Context) {
	var sum models.InputWithdraw

	id, err := strconv.Atoi(ctx.Param("account_id"))

	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&sum); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.WithdrawProducer(id, sum)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) transfer(ctx *gin.Context) {
	var sum models.InputTransfer

	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&sum); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.TransferProducer(id, sum)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) transactions(ctx *gin.Context) {

}
