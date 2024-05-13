package flutterwave_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"usepolymer.co/background/logger"
	"usepolymer.co/background/network"
)


var LocalPaymentProcessor *FlutterwavePaymentProcessor = &FlutterwavePaymentProcessor{}


type FlutterwavePaymentProcessor struct {
	Network *network.NetworkController
	AuthToken string
}

func (fpp *FlutterwavePaymentProcessor) InitialisePaymentProcessor() {
	LocalPaymentProcessor.Network =  &network.NetworkController{
		BaseUrl: os.Getenv("FLUTTERWAVE_BASE_URL"),
	}
	LocalPaymentProcessor.AuthToken =  os.Getenv("FLUTTERWAVE_ACCESS_TOKEN")
}

func (fpp *FlutterwavePaymentProcessor) GetTransactionDetail(id int64) (*string) {
	response, statusCode, err := fpp.Network.Get(fmt.Sprintf("/transfers/%d", id), &map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", fpp.AuthToken),
	}, nil)
	if err != nil {
		logger.Error(errors.New("an error occured while retrieving pending transaction on flutterwave"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}

	var flwResponse TransferDetailResponse
	err = json.Unmarshal(*response, &flwResponse)

	if err != nil {
		logger.Error(errors.New("an error occured while unmarshalling pending transaction payload from flutterwave"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}

	if *statusCode != 200 {
		err = errors.New("retrieving pending transaction on flutterwave")
		logger.Error(err, logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "body",
			Data: flwResponse,
		})
		return nil
	}
	return &flwResponse.Data.Status
}
