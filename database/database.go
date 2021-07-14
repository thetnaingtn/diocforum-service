package database

import (
	"context"
	_ "diocforum/cert"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func init() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(os.Getenv("DATABASE_URL"), filepath.Join("ca-certificate.crt"))))
	Client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(os.Getenv("DATABASE_URL"), filepath.Join("cert/ca-certificate.crt"))))
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(os.Getenv("DATABASE_URL"))
	log.Println(Client.ListDatabaseNames(ctx, bson.M{}))
}
