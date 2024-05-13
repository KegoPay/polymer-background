package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"usepolymer.co/background/controllers"
	"usepolymer.co/background/logger"
)

var Router *gin.RouterGroup

func WalletRouter() {
	walletRouter := Router.Group("/wallet")
	
	walletRouter.POST("/request-statement", func(ctx *gin.Context) {
		var body controllers.RequestAccountStatementDTO
		if err := ctx.ShouldBindJSON(&body); err != nil {
			logger.Error(errors.New("error binding payload"), logger.LoggerOptions{
				Key: "error",
				Data: err,
			})
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		err := controllers.RequestAccountStatement(&body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	})
}