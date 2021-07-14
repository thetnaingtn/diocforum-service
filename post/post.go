package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Body      string             `json:"body,omitempty" bson:"body,omitempty"`
	AvatarURL string             `json:"avatarURL,omitempty" bson:"avatarURL,omitempty"`
	UserName  string             `json:"userName,omitempty" bson:"userName,omitempty"`
	UserId    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	ThreadId  primitive.ObjectID `json:"threadId,omitempty" bson:"threadId,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
