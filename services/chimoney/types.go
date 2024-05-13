package chimoney_service

import "time"

type ChimoneyTransactionDetail struct {
	Data ChimoneyTransaction `json:"data"`
}

type ChimoneyTransaction struct {
	Status 				string `json:"status"`
	DeliveryStatus  	string `json:"deliveryStatus"`
	RecievedAt		  	time.Time `json:"redeemDate"`
	Payout  			ChimoneyPayoutError `json:"payout"`
}

type ChimoneyPayoutError struct {
	Error *string	`json:"error"`
}