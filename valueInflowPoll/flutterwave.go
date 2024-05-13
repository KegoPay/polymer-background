package valueinflowpoll

import (
	// "errors"
	// "fmt"

	// "go.mongodb.org/mongo-driver/mongo/options"
	// "usepolymer.co/background/logger"
	// "usepolymer.co/background/repository"
)

var BATCH_LIMIT int64 = 500

func pollFlutterwave() {
	// walletRepository := repository.WalletRepo()
	// var lastID string

	// func () {
	// 	wallets, err := walletRepository.FindMany(map[string]interface{}{
	// 		"_id": map[string]interface{}{"$lt": lastID},
	// 	}, options.Find().SetSort(map[string]interface{}{
	// 		"_id": -1,
	// 	}), &options.FindOptions{
	// 		Limit: &BATCH_LIMIT,
	// 	})
	// 	if err != nil {
	// 		logger.Error(fmt.Errorf("error fetching batch"), logger.LoggerOptions{
	// 			Key: "lastID",
	// 			Data: lastID,
	// 		})
	// 	}
	// 	for _, wallet := range *wallets {
	// 		go func() {
	// 			wallet.
	// 		}()
	// 	}
	// }()
}