package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	AvatarURL    string             `json:"avatarURL,omitempty" bson:"avatarURL,omitempty"`
	ThreadId     primitive.ObjectID `json:"threadId,omitempty" bson:"threadId,omitempty"`
	RegisteredAt time.Time          `json:"registeredAt,omitempty" bson:"registeredAt,omitempty"`
}

type Identity struct {
	UserId string `json:"userId,omitempty"`
}
