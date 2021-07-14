package user

import (
	"context"
	"diocforum/database"
	"diocforum/util"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUserExist error
var db = database.Client.Database("forum")
var collection = db.Collection("user")

func CreateUser(ctx context.Context, user User) (primitive.ObjectID, error) {
	var result User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&result)
	if err != mongo.ErrNoDocuments {
		ErrUserExist = fmt.Errorf("User %s already exist", result.Name)
		return result.Id, ErrUserExist
	}

	user.RegisteredAt = time.Now()
	user.Password = util.Encrypt(user.Password)
	insertedResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertedResult.InsertedID.(primitive.ObjectID), nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
