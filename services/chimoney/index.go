package chimoney_service

import (
	"encoding/json"
	"errors"
	"os"

	"usepolymer.co/background/logger"
	"usepolymer.co/background/network"
)



var InternationalPaymentProcessor *ChimoneyPaymentProcessor = &ChimoneyPaymentProcessor{}

type ChimoneyPaymentProcessor struct {
	Network *network.NetworkController
	AuthToken string
}

func (chimoneyPP *ChimoneyPaymentProcessor) InitialisePaymentProcessor() {
	InternationalPaymentProcessor.Network = &network.NetworkController{
		BaseUrl: os.Getenv("CHIMONEY_BASE_URL"),
	}
	InternationalPaymentProcessor.AuthToken = os.Getenv("CHIMONEY_ACCESS_TOKEN")
}


func (chimoneyPP *ChimoneyPaymentProcessor) GetTransactionDetail(trxID string) (*ChimoneyTransaction){
	response, statusCode, err := chimoneyPP.Network.Post("/payment/verify", &map[string]string{
		"X-API-KEY": chimoneyPP.AuthToken,
	}, map[string]interface{}{
		"id": trxID,
	}, nil)

	if err != nil {
		logger.Error(errors.New("an error occured while polling transaction detail on chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "trxID",
			Data: trxID,
		})
		return nil
	}
	var chimoneyResponse ChimoneyTransactionDetail
	err = json.Unmarshal(*response, &chimoneyResponse)
	if err != nil {
		logger.Error(errors.New("an error occured while parsing transaction detail from chimoney"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	if *statusCode != 200 {
		err = errors.New("an error occured while polling transaction detail on chimoney")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: chimoneyResponse,
		})
		return nil
	}
	return &chimoneyResponse.Data
}
