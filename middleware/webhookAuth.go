package middleware

import (
	"github.com/gin-gonic/gin"
	"usepolymer.co/background/logger"
)

func ChimoneyWebHookAuth(ctx *gin.Context) {
	logger.Info("this should authenticate the webhook request")
}