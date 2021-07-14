package thread

import (
	"context"
	"diocforum/database"
	"diocforum/post"
	"diocforum/user"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db = database.Client.Database("forum")
var collection = db.Collection("thread")
var postCollection = db.Collection("post")
var userCollection = db.Collection("user")

func createThread(ctx context.Context, thread Thread) (*mongo.InsertOneResult, error) {
	thread.CreatedAt = time.Now()
	result, err := collection.InsertOne(ctx, thread)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getThread(ctx context.Context, id primitive.ObjectID) (*Thread, error) {
	var thread Thread

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thread)
	if err != nil {
		return nil, err
	}

	cursor, err := postCollection.Find(ctx, bson.M{"threadId": id})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	//need refactor
	for cursor.Next(ctx) {
		var postUser user.User
		var post post.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		userCollection.FindOne(ctx, bson.M{"_id": post.UserId}).Decode(&postUser)
		post.AvatarURL = postUser.AvatarURL
		post.UserName = postUser.Name
		thread.Posts = append(thread.Posts, post)
	}

	return &thread, nil

}

func getAllThreads(ctx context.Context) ([]Thread, error) {
	var threads []Thread
	findOption := options.Find().SetSort(bson.D{{"createdAt", -1}})

	threadCursor, err := collection.Find(ctx, bson.M{}, findOption)
	defer threadCursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	//need refactor
	for threadCursor.Next(ctx) {
		var thread Thread
		var threadUser user.User
		var posts []post.Post
		if err := threadCursor.Decode(&thread); err != nil {
			return nil, err
		}
		userCollection.FindOne(ctx, bson.M{"_id": thread.UserId}).Decode(&threadUser)
		thread.UserName = threadUser.Name
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
		postCursor, _ := db.Collection("post").Find(ctx, bson.M{"threadId": thread.Id}, findOptions)
		for postCursor.Next(ctx) {
			var post post.Post
			var postUser user.User
			if err := postCursor.Decode(&post); err != nil {
				return nil, err
			}
			userCollection.FindOne(ctx, bson.M{"_id": post.UserId}).Decode(&postUser)
			post.AvatarURL = postUser.AvatarURL
			posts = append(posts, post)
		}
		thread.Posts = posts
		threads = append(threads, thread)

	}

	return threads, nil

}

func deleteThread(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	_, err = postCollection.DeleteMany(ctx, bson.M{"threadId": id})
	if err != nil {
		return nil, fmt.Errorf("Error occour when posts are deleted by thread id: %s", err.Error())
	}

	return result, nil
}

func updateThread(ctx context.Context, id primitive.ObjectID, thread Thread) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: thread}},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
