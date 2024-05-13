package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"usepolymer.co/background/constants"
	"usepolymer.co/background/database/connection"
	"usepolymer.co/background/logger"
	chimoney_service "usepolymer.co/background/services/chimoney"
	flutterwave_service "usepolymer.co/background/services/flutterwave"
	valueinflowpoll "usepolymer.co/background/valueInflowPoll"
	valueoutflowpoll "usepolymer.co/background/valueOutFlowPoll"
)

func main() {
	godotenv.Load()

	logger.InitializeLogger()
	connection.ConnectToDatabase()
	constants.InitialisePollingIntervals()
	chimoney_service.InternationalPaymentProcessor.InitialisePaymentProcessor()
	flutterwave_service.LocalPaymentProcessor.InitialisePaymentProcessor()
	valueoutflowpoll.PollForValueOutflow()
	valueinflowpoll.PollForValueInflow()

	server := gin.Default()

    v1 := server.Group("/api/v1")
    Router = v1
    WalletRouter()

	gin_mode := os.Getenv("GIN_MODE")
	port := os.Getenv("PORT")
	if gin_mode == "debug" || gin_mode == "release"{
		logger.Info(fmt.Sprintf("Server starting on PORT %s", port))
		server.Run(fmt.Sprintf(":%s", port))
	} else {
		panic(fmt.Sprintf("invalid gin mode used - %s", gin_mode))
	}
}
