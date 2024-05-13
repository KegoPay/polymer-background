package database

import "usepolymer.co/background/database/connection"

func SetUpDatabase(){
	connection.ConnectToDatabase()
}

type BaseModel interface {
	ParseModel() any
}
