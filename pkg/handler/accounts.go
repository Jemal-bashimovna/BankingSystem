package handler

import (
	"bankingsystem/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createAccount(ctx *gin.Context) {
	var inputAccount models.CreateAccount

	if err := ctx.BindJSON(&inputAccount); err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateAccount(inputAccount)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.CreateAccResponse{
		Status: "Successfully",
		Id:     id,
	})
}

func (h *Handler) deleteAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Accounts.DeleteAccount(id)
	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.DeleteAccResponse{
		Status: "Successfully",
	})

}

func (h *Handler) getAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		ErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.service.Accounts.GetAccountById(id)

	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, models.GetAccResponse{
		Account: account,
	})
}

func (h *Handler) getAllAccounts(ctx *gin.Context) {
	accounts, err := h.service.Accounts.GetAccounts()

	if err != nil {
		ErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.GetAllAccResponse{
		Accounts: accounts,
	})
}
