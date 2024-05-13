package repository

import (
	"sync"

	"usepolymer.co/background/database/connection/datastore"
	"usepolymer.co/background/database/repository/mongo"
	"usepolymer.co/background/models"
)


var walletOnce = sync.Once{}

var walletRepository mongo.MongoRepository[models.Wallet]

func WalletRepo() *mongo.MongoRepository[models.Wallet] {
	walletOnce.Do(func() {
		walletRepository = mongo.MongoRepository[models.Wallet]{Model: datastore.WalletModel}
	})
	return &walletRepository
}
