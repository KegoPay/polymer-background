package repository

import (
	"sync"

	"usepolymer.co/background/database/connection/datastore"
	"usepolymer.co/background/database/repository/mongo"
	"usepolymer.co/background/models"
)


var transactionOnce = sync.Once{}

var transactionRepository mongo.MongoRepository[models.Transaction]

func TransactionRepo() *mongo.MongoRepository[models.Transaction] {
	transactionOnce.Do(func() {
		transactionRepository = mongo.MongoRepository[models.Transaction]{Model: datastore.TransactionModel}
	})
	return &transactionRepository
}
