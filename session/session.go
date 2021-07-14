package session

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	Id     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Email  string             `json:"email,omitempty" bson:"email,omitempty"`
}
