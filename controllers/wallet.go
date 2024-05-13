package controllers

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/messaging/emails"
	"usepolymer.co/background/repository"
	toexcel "usepolymer.co/background/services/toExcel"
	"usepolymer.co/background/utils"
)

func RequestAccountStatement(body *RequestAccountStatementDTO) error {
	transactionRepo := repository.TransactionRepo()
	startDate, err := time.Parse("2006-01-02", body.Start) 
	if err != nil {
		return err
	}
	endDate, err := time.Parse("2006-01-02", body.End) 
	if err != nil {
		return err
	}
	transactions, err := transactionRepo.FindMany(map[string]interface{}{
		"walletID": body.WalletID,
		"createdAt": map[string]any{
			"$gte": startDate,
			"$lte": endDate,
		},
	}, options.Find().SetProjection(map[string]any{
		"createdAt": 1,
		"transactionReference": 1,
		"description": 1,
		"amount": 1,
		"currency": 1,
		"intent": 1,
		"transactionRecepient": 1,
	}))
	if err != nil {
		logger.Error(errors.New("error fetching transaction data for account sttaement generation"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	fileName, err := toexcel.TransactionToExcel(transactions)
	if err != nil {
		logger.Error(errors.New("error converting transaction to excel for account sttaement generation"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	projectRoot, _ := os.Getwd()
	filePath := filepath.Join(projectRoot, *fileName)
	file, err := utils.LoadFile(filePath)
	if err != nil {
		logger.Error(errors.New("error loading account statement file"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
	}
	utils.DeleteFile(filePath)
	emails.EmailService.SendEmail(body.Email, "Polymer Account Statement", "statement_generated", nil, &emails.ResendAttachment{
		Name: *fileName,
		Content: file,
	})
	return nil
}