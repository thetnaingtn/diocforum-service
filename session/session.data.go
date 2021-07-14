package session

import (
	"context"
	"diocforum/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db = database.Client.Database("forum")
var collection = db.Collection("session")

func CreateSession(ctx context.Context, session Session) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetSessionByUserId(ctx context.Context, id string) (*Session, error) {
	var session Session
	objectId, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, bson.M{"userId": objectId}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func DeleteUserSession(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteMany(ctx, bson.M{"userId": id})
	if err != nil {
		return nil, err
	}

	return result, nil
}
