package valueoutflowpoll

import (
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/background/database/repository/mongo"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/models"

	"usepolymer.co/background/repository"
	flutterwave_service "usepolymer.co/background/services/flutterwave"
	wallet_service "usepolymer.co/background/services/wallet"
)

var FLUTTERWAVE_BATCH_LIMIT int64 = 500

func pollFlutterwave() {
	var wg sync.WaitGroup
	transactionRepository := repository.TransactionRepo()
	var lastID string
	pendingTrxCount, err := transactionRepository.CountDocs(map[string]interface{}{
		"intent": models.LocalDebit,
		"status": models.TransactionPending,
	})
	if err != nil {
		logger.Error(errors.New("error counting pending domestic outflowing transactions"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return
	}
	if pendingTrxCount == 0 {
		logger.Info("no pending domestic outflowing transactions to process")
		return
	}
	processFlutterwaveBatch(transactionRepository, lastID, &wg, pendingTrxCount)
}

func processFlutterwaveBatch(transactionRepository *mongo.MongoRepository[models.Transaction], lastID string, wg *sync.WaitGroup, pendingTrxCount int64) {
	transactions, err := transactionRepository.FindMany(map[string]interface{}{
		"_id": func() map[string]interface{} {
			if lastID == "" {
				return map[string]interface{}{"$gt": ""}
			}
			return    map[string]interface{}{"$lt": lastID}
		}(),
		"intent": models.LocalDebit,
		"status": models.TransactionPending,
	}, options.Find().SetSort(map[string]interface{}{
		"_id": -1,
	}), &options.FindOptions{
		Limit: &FLUTTERWAVE_BATCH_LIMIT,
	})
	if err != nil {
		logger.Error(errors.New("error fetching batch"), logger.LoggerOptions{
			Key: "lastID",
			Data: lastID,
		})
	}
	if transactions == nil || len(*transactions) == 0 {
		logger.Info("no pending domestic outflowing transactions to process")
		return
	}
	lastID = (*transactions)[len(*transactions) - 1].ID
	for _, trx := range *transactions {
		wg.Add(1)
		go func(t models.Transaction) {
			status := flutterwave_service.LocalPaymentProcessor.GetTransactionDetail(t.MetaData["paymentid"].(int64))
			if status == nil {
				logger.Info("no flutterwave transaction for polymer transaction", logger.LoggerOptions{
					Key: "id",
					Data: t.ID,
				})
				return
			}
			switch *status {
				case "FAILED": 
					wallet_service.ReverseLockFunds(t.WalletID, t.TransactionReference)
					affected, err := transactionRepository.UpdatePartialByID(t.ID, map[string]any{
						"status": models.TransactionFailed,
					})
					if err != nil {
						logger.Error(errors.New("error updating flutterwave transaction after refund"), logger.LoggerOptions{
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
				case "SUCCESSFUL":
					wallet_service.RemoveLockFunds(t.WalletID, t.TransactionReference)
					affected, err := transactionRepository.UpdatePartialByID(t.ID, map[string]any{
						"status": models.TransactionCompleted,
					})
					if err != nil {
						logger.Error(errors.New("error updating flutterwave transaction after success"), logger.LoggerOptions{
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
	pendingTrxCount -= FLUTTERWAVE_BATCH_LIMIT
	if pendingTrxCount <= 0 {
		return
	}else {
		processFlutterwaveBatch(transactionRepository, lastID, wg, pendingTrxCount)
	}
}