package flutterwave_service

type TransferDetailResponse struct {
	Data	*TransferDetail		`json:"data"`
}

type TransferDetail struct {
	Status 	string	`json:"status"`
}