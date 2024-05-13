package repository

import (
	"sync"

	"usepolymer.co/background/database/connection/datastore"
	"usepolymer.co/background/database/repository/mongo"
	"usepolymer.co/background/models"
)


var userOnce = sync.Once{}

var userRepository mongo.MongoRepository[models.User]

func UserRepo() *mongo.MongoRepository[models.User] {
	userOnce.Do(func() {
		userRepository = mongo.MongoRepository[models.User]{Model: datastore.UserModel}
	})
	return &userRepository
}
