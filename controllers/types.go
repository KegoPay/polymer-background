package controllers

type RequestAccountStatementDTO struct {
	Start 		string  `json:"start"`
	Email 		string  `json:"email"`
	End	  	 	string  `json:"end"`
	WalletID	string  `json:"walletID"`
}