package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorMessage struct {
	Message string `json:"message"`
}

func ErrorMessage(ctx *gin.Context, httpStatus int, message string) {
	logrus.Errorf(message)
	ctx.AbortWithStatusJSON(httpStatus, errorMessage{Message: message})
}
