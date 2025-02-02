package models

import (
	"fmt"
	"time"

	"usepolymer.co/background/utils"
)


type PhoneNumber struct	{
	ISOCode 	 string `bson:"isoCode" json:"isoCode"` // Two-letter country code (ISO 3166-1 alpha-2)
	LocalNumber  string `bson:"localNumber" json:"localNumber"`
	Prefix		 string `bson:"prefix" json:"prefix"`
	IsVerified	 bool   `bson:"isVerified" json:"isVerified"`
	WhatsApp	 bool   `bson:"whatsapp" json:"whatsapp"`
}

func (pn *PhoneNumber) ParsePhoneNumber() string {
	return fmt.Sprintf("+%s%s", pn.Prefix, pn.LocalNumber)
}

type NotificationOptions struct {
	PushNotification bool `bson:"pushNotification" json:"pushNotification"`
	Emails 			 bool `bson:"emails" json:"emails"`
}

type NextOfKin struct {
	FirstName 		string `bson:"firstName" json:"firstName"`
	LastName 		string `bson:"lastName" json:"lastName"`
	Relationship 	string `bson:"relationship" json:"relationship"`
}

type Address struct {
	FullAddress 	*string 	`bson:"fullAddress" json:"fullAddress"`
	Verified 		bool 	`bson:"verified" json:"verified"`
	State 			*string 	`bson:"state" json:"state"`
	LGA 			*string 	`bson:"lga" json:"lga"`
	Street 			*string 	`bson:"street" json:"street"`
}


type User struct {
	FirstName         				string       			`bson:"firstName" json:"firstName"`
	LastName          				string       			`bson:"lastName" json:"lastName"`
	MiddleName          			*string      			`bson:"middleName" json:"middleName"`
	Email            				string       			`bson:"email" json:"email" validate:"required,email"`
	Phone             				*PhoneNumber  			`bson:"phone" json:"phone"`
	Password          				string       			`bson:"password" json:"-" validate:"required,password"`
	TransactionPin    				string       			`bson:"transactionPin" json:"-"`
	UserAgent        				string    	 			`bson:"userAgent" json:"-" validate:"required"`
	DeviceID          				string       			`bson:"deviceID" json:"-" validate:"required"`
	PushNotificationToken	        string       			`bson:"pushNotificationToken" json:"-" validate:"required"`
	AppVersion          			string       			`bson:"appVersion" json:"-"`
	WalletID  						string    				`bson:"walletID" json:"walletID"`
	KYCCompleted   					bool         			`bson:"kycCompleted" json:"kycCompleted"`
	HasBusiness   					bool         			`bson:"hasBusiness" json:"hasBusiness"`
	EmailVerified     				bool         			`bson:"emailVerified" json:"emailVerified"`
	AccountRestricted 				bool         			`bson:"accountRestricted" json:"accountRestricted"`
	AccountLocked	 				bool         			`bson:"accountLocked" json:"accountLocked"`
	Deactivated 					bool         			`bson:"deactivated" json:"deactivated"`
	BVN		  		  				string 	  	 			`bson:"bvn" json:"-"`
	NIN		  		  				string 	  	 			`bson:"nin" json:"-"`
	Gender		  		  			string 	  	 			`bson:"gender" json:"gender"`
	Address							*Address				`bson:"address" json:"address"`
	DOB		  		  				string 	  	 			`bson:"dob" json:"dob"`
	WatchListed		  		  		bool 	  	 			`bson:"watchListed" json:"-"`
	NINLinked		  		  		bool 	  	 			`bson:"ninLinked" json:"ninLinked"`
	Nationality		  		  		string 	  	 			`bson:"nationality" json:"nationality"`
	ProfileImage		  		  	string 	  	 			`bson:"profileImage" json:"profileImage"`
	Tag		  		  				string 	  			 	`bson:"tag" json:"tag"`
	Longitude		  		  		float64 	  			`bson:"longitude" json:"longitude"`
	Latitude		  		  		float64 	  			`bson:"latitude" json:"latitude"`
	Tier		  		  			uint 	  			 	`bson:"tier" json:"tier"`
	NotificationOptions		  		NotificationOptions  	`bson:"notificationOptions" json:"notificationOptions"`
	NextOfKin		  				NextOfKin  				`bson:"nextOfKin" json:"nextOfKin"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (user User) ParseModel() any {
	now := time.Now()
	if user.ID == "" {
		user.CreatedAt = now
		user.ID = utils.GenerateUUIDString()
	}
	user.UpdatedAt = now
	return &user
}
