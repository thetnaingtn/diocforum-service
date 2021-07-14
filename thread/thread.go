package thread

import (
	"diocforum/post"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thread struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Topic     string             `json:"topic,omitempty" bson:"topic,omitempty"`
	UserId    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	UserName  string             `json:"userName,omitempty" bson:"userName,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Posts     []post.Post        `json:"posts,omitempty" bson:"posts,omitempty"`
}
