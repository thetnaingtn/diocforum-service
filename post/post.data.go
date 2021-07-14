package post

import (
	"context"
	"diocforum/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db = database.Client.Database("forum")
var collection = db.Collection("post")

func getPost(ctx context.Context, id primitive.ObjectID) (*Post, error) {
	var post Post
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func createPost(ctx context.Context, post Post) (*mongo.InsertOneResult, error) {
	post.CreatedAt = time.Now()
	result, err := collection.InsertOne(ctx, post)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func updatePost(ctx context.Context, post Post, id primitive.ObjectID) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: post}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func deletePost(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return result, nil
}
