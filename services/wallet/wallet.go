package wallet_service

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/models"
	"usepolymer.co/background/repository"
)

func ReverseLockFunds(walletID string, lockedFundsReference string) error {
	walletRepository := repository.WalletRepo()
	wallet, err := walletRepository.FindByID(walletID)
	if err != nil {
		logger.Error(errors.New("could not reverse lock funds"), logger.LoggerOptions{
			Key: "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	if wallet == nil {
		logger.Error(errors.New("wallet not found"), logger.LoggerOptions{
			Key: "walletID",
			Data: wallet,
		})
		return err
	}
	var lockedFunds []models.LockedFunds
	for _, lf := range wallet.LockedFundsLog {
		if lf.LockedFundsID == lockedFundsReference {
			lockedFunds = append(lockedFunds, lf)
		}
	}
	var totalAmount uint64 = 0
	for _, lockedFund := range lockedFunds {
		totalAmount += lockedFund.Amount
	}
	affected, err := walletRepository.UpdateWithOperator(context.TODO(), map[string]interface{}{
		"_id": walletID,
	}, map[string]any{
		"$pull": map[string]any {
			"lockedFundsLog": map[string]any{
				"$in": lockedFunds,
			},
		},
		"$inc": map[string]any {
			"balance": totalAmount,
		},
	})

	if !affected {
		logger.Error(errors.New("could not reverse lock funds"), logger.LoggerOptions{
			Key: "walletID",
			Data: walletID,
		})
		return err
	}
	return nil
}

func RemoveLockFunds(walletID string, lockedFundsReference string) error {
	walletRepository := repository.WalletRepo()
	wallet, err := walletRepository.FindByID(walletID, options.FindOne().SetProjection(map[string]any{
		"lockedFundsLog": 1,
	}))
	if err != nil {
		logger.Error(errors.New("could not remove lock funds"), logger.LoggerOptions{
			Key: "walletID",
			Data: walletID,
		}, logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	if wallet == nil {
		logger.Error(errors.New("wallet not found"), logger.LoggerOptions{
			Key: "walletID",
			Data: wallet,
		})
		return err
	}
	var polymer_vat models.LockedFunds
	var polymer_fee models.LockedFunds
	var lockedFunds []models.LockedFunds
	for _, lf := range wallet.LockedFundsLog {
		if lf.LockedFundsID == lockedFundsReference {
			lockedFunds = append(lockedFunds, lf)
			if lf.Reason == models.PolymerFee {
				polymer_fee = lf
			}
			if lf.Reason == models.PolymerVAT {
				polymer_vat = lf
			}
		}
	}
	var totalAmount uint64 = 0
	for _, lockedFund := range lockedFunds {
		totalAmount += lockedFund.Amount
	}
	walletRepository.StartTransaction(func(sc mongo.Session, c context.Context) error {
		
		affected, e := walletRepository.UpdateWithOperator(c, map[string]interface{}{
			"_id": walletID,
		}, map[string]any{
			"$pull": map[string]any {
				"lockedFundsLog": map[string]any{
					"$in": lockedFunds,
				},
			},
			"$inc": map[string]any {
				"ledgerBalance": -int64(totalAmount),
			},
		})

		if !affected {
			logger.Error(errors.New("could not remove lock funds"), logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			return err
		}
		if e != nil {
			logger.Error(errors.New("could not remove lock funds"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "lockedFunds",
				Data: lockedFunds,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		if !affected {
			logger.Error(errors.New("could not update wallet"), logger.LoggerOptions{
				Key: "lockedFunds",
				Data: lockedFunds,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}

		affected, err = walletRepository.UpdateWithOperator(c, map[string]interface{}{
			"userID": os.Getenv("POLYMER_WALLET_USER_ID"),
			"businessName": "Polymer Fee Wallet",
		}, map[string]any{
			"$inc": map[string]any {
				"ledgerBalance": polymer_fee.Amount,
				"balance": polymer_fee.Amount,
			},
		})

		if e != nil || !affected {
			logger.Error(errors.New("could not create transaction entry"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "polymer vat lf",
				Data: polymer_vat,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			}, logger.LoggerOptions{
				Key: "affected",
				Data: affected,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}

		affected, err = walletRepository.UpdateWithOperator(c, map[string]interface{}{
			"userID": os.Getenv("POLYMER_WALLET_USER_ID"),
			"businessName": "Polymer VAT Wallet",
		}, map[string]any{
			"$inc": map[string]any {
				"ledgerBalance": polymer_vat.Amount,
				"balance": polymer_vat.Amount,
			},
		})

		if e != nil || !affected {
			logger.Error(errors.New("could not create transaction entry"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "polymer vat lf",
				Data: polymer_vat,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			}, logger.LoggerOptions{
				Key: "affected",
				Data: affected,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		(sc).CommitTransaction(c)
		return nil
	})
	if err != nil {
		logger.Error(errors.New("an error occured while removing locked funds and crediting default wallets"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
	}
	return err
}

func CreditWallet(walletID string, amount uint64, intent models.TransactionIntent, transactionPayload *models.Transaction) error {
	var err error
	walletRepository := repository.WalletRepo()
	transactionRepository := repository.TransactionRepo()
	walletRepository.StartTransaction(func(sc mongo.Session, c context.Context) error {
		affected, e := walletRepository.UpdateManyWithOperator(c, map[string]interface{}{
			"_id": walletID,
		}, map[string]any{
			"$inc": map[string]any {
				"balance": amount,
				"ledgerBalance": amount,
			},
		})
		if e != nil {
			logger.Error(errors.New("could not credit account"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "transaction",
				Data: transactionPayload,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		if affected != 1 {
			logger.Error(errors.New("could not credit account or multiple accounts credited"), logger.LoggerOptions{
				Key: "transaction",
				Data: transactionPayload,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}

		_, e = transactionRepository.CreateOne(c, *transactionPayload)
		if e != nil {
			logger.Error(errors.New("could not create transaction entry"), logger.LoggerOptions{
				Key: "error",
				Data: e,
			}, logger.LoggerOptions{
				Key: "transaction",
				Data: transactionPayload,
			},logger.LoggerOptions{
				Key: "walletID",
				Data: walletID,
			})
			err = e
			(sc).AbortTransaction(c)
			return e
		}
		(sc).CommitTransaction(c)
		return nil
	})
	return err
}
