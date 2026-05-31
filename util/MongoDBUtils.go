package util

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbConnect() (*mongo.Client, context.Context) {
	ctx := context.Background()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"code": 1},
		Options: options.Index().SetUnique(true),
	}

	client.Database("testing").Collection("hashKeyPairs").Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		log.Fatal()
	}
	return client, ctx
}
