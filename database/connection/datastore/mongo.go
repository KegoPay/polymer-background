package datastore

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/background/logger"
)
 
var (
	UserModel *mongo.Collection
	TransactionModel *mongo.Collection
	WalletModel *mongo.Collection
)

func connectMongo() *context.CancelFunc {
	url := os.Getenv("DB_URL")

	if url == "" {
		logger.Error(errors.New("set mongo url"))
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	clientOpts := options.Client().ApplyURI(url)
	clientOpts.SetMinPoolSize(5)
	clientOpts.SetMaxPoolSize(10)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		logger.Warning("an error occured while starting the database", logger.LoggerOptions{Key: "error", Data: err})
		return &cancel
	}

	db := client.Database(os.Getenv("DB_NAME"))
	setUpIndexes(ctx, db)

	logger.Info("connected to mongodb successfully")
	return &cancel
}

// Set up the indexes for the database
func setUpIndexes(ctx context.Context, db *mongo.Database) {
	UserModel = db.Collection("Users")
	WalletModel = db.Collection("Wallets")
	TransactionModel = db.Collection("Transactions")
}
