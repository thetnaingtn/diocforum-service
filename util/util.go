package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Encrypt(plaintext string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(plaintext))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func GetObjectIdFromUrl(url, spliter string) (primitive.ObjectID, error) {
	hexString := strings.Split(url, spliter)[1]
	if len(hexString) < 1 {
		return primitive.NilObjectID, fmt.Errorf("Hexstring is too short")
	}
	objectId, err := primitive.ObjectIDFromHex(hexString[1:])
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objectId, nil
}
