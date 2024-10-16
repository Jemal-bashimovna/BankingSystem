package handler

import (
	"bankingsystem/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) deposit(ctx *gin.Context) {

	var input models.InputDeposit

	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&input); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.DepositProducer(id, input)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.TransactionResponse{
		Message: fmt.Sprintf("Deposit to account: %d successfully", input.AccountId),
	})
}

func (h *Handler) withdraw(ctx *gin.Context) {
	var input models.InputWithdraw

	id, err := strconv.Atoi(ctx.Param("account_id"))

	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&input); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.WithdrawProducer(id, input)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.TransactionResponse{
		Message: fmt.Sprintf("Withdrawing money (%.2f) successfully from account: %d", input.WithDrawSum, input.AccountId),
	})
}

func (h *Handler) transfer(ctx *gin.Context) {
	var input models.InputTransfer

	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.BindJSON(&input); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Transactions.TransferProducer(id, input)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.TransactionResponse{
		Message: fmt.Sprintf("The transfer %.2f from: %d to %d was successfully", input.TransferSum, input.AccountId, input.TargetId),
	})
}

func (h *Handler) transactions(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}
	transactions, err := h.service.Transactions.GetAll(id)

	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, models.GetTransactionsResponse{
		Transactions: transactions,
	})
}
