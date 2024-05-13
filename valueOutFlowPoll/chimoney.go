package valueoutflowpoll

import (
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/background/database/repository/mongo"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/models"

	"usepolymer.co/background/repository"
	chimoney_service "usepolymer.co/background/services/chimoney"
	wallet_service "usepolymer.co/background/services/wallet"
)

var CHIMONEY_BATCH_LIMIT int64 = 500

func pollChimoney() {
	var wg sync.WaitGroup
	transactionRepository := repository.TransactionRepo()
	var lastID string
	pendingTrxCount, err := transactionRepository.CountDocs(map[string]interface{}{
		"intent": models.InternationalDebit,
		"status": models.TransactionPending,
	})
	if err != nil {
		logger.Error(errors.New("error counting pending international transactions"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}
	if pendingTrxCount == 0 {
		logger.Info("no pending international transactions to process")
		return
	}
	processChimoneyBatch(transactionRepository, lastID, &wg, pendingTrxCount)
}

func processChimoneyBatch(transactionRepository *mongo.MongoRepository[models.Transaction], lastID string, wg *sync.WaitGroup, pendingTrxCount int64) {
	transactions, err := transactionRepository.FindMany(map[string]interface{}{
		"_id": func() map[string]interface{} {
			if lastID == "" {
				return map[string]interface{}{"$gt": ""}
			}
			return    map[string]interface{}{"$lt": lastID}
		}(),
		"intent": models.InternationalDebit,
		"status": models.TransactionPending,
	}, options.Find().SetSort(map[string]interface{}{
		"_id": -1,
	}), &options.FindOptions{
		Limit: &CHIMONEY_BATCH_LIMIT,
	})
	if err != nil {
		logger.Error(errors.New("error fetching batch"), logger.LoggerOptions{
			Key: "lastID",
			Data: lastID,
		})
	}
	if transactions == nil || len(*transactions) == 0 {
		logger.Info("no pending international transactions to process")
		return
	}
	lastID = (*transactions)[len(*transactions) - 1].ID
	for _, trx := range *transactions {
		wg.Add(1)
		go func(t models.Transaction) {
			transaction := chimoney_service.InternationalPaymentProcessor.GetTransactionDetail(t.MetaData["issueid"].(string))
			if transaction == nil {
				logger.Info("no chimoney transaction for polymer transaction", logger.LoggerOptions{
					Key: "id",
					Data: t.ID,
				})
				return
			}
			var message string
			if transaction.Payout.Error != nil {
				if *transaction.Payout.Error == "Invalid account_bank and destination_branch_code combination passed." {
					message = "Payment failed because there was a bank account and branch code mixup"
				}else {
					message = "Payment failed due to an unknown reason. We are currently investigating the cause."
				}
			}else {
				message = "Payment was successful!🥳🎉🤑"
			}
			switch transaction.Status {
				case "refunded": 
					wallet_service.ReverseLockFunds(t.WalletID, t.TransactionReference)
					affected, err := transactionRepository.UpdatePartialByID(t.ID, map[string]any{
						"status": models.TransactionFailed,
						"message": message,
					})
					if err != nil {
						logger.Error(errors.New("error updating chimoney transaction after refund"), logger.LoggerOptions{
							Key: "error",
							Data: err,
						}, logger.LoggerOptions{
							Key: "id",
							Data: t.ID,
						})
					}
					if affected != 1 {
						logger.Error(errors.New("transaction refund update failed to affect 1 transaction"), logger.LoggerOptions{
							Key: "id",
							Data: t.ID,
						})
					}
				case "redeemed":
					wallet_service.RemoveLockFunds(t.WalletID, t.TransactionReference)
					affected, err := transactionRepository.UpdatePartialByID(t.ID, map[string]any{
						"status": models.TransactionCompleted,
					})
					if err != nil {
						logger.Error(errors.New("error updating chimoney transaction after success"), logger.LoggerOptions{
							Key: "error",
							Data: err,
						}, logger.LoggerOptions{
							Key: "id",
							Data: t.ID,
						})
					}
					if affected != 1 {
						logger.Error(errors.New("transaction success update failed to affect 1 transaction"), logger.LoggerOptions{
							Key: "id",
							Data: t.ID,
						})
					}
			}
			wg.Done()
		}(trx)
	}
	wg.Wait()
	pendingTrxCount -= CHIMONEY_BATCH_LIMIT
	if pendingTrxCount <= 0 {
		return
	}else {
		processChimoneyBatch(transactionRepository, lastID, wg, pendingTrxCount)
	}
}