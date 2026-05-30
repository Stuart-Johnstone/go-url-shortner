package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type kvPair struct {
	shortened string
	original  string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET/{id}", handleRedirect)
	http.ListenAndServe(":8080", mux)

	// var userInput string
	// fmt.Scanln(&userInput)
	// userKv := getHash(userInput)
	//
	// client, ctx := dbConnect()
	//
	// writeToDb(userKv, client, ctx)
	// fmt.Scanln(&userInput)
	// readFromDb(userInput, client, ctx)
}

func dbConnect() (*mongo.Client, context.Context) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

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

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("id")

	doc := bson.M{
		"code": path,
	}

	client, ctx := dbConnect()

	collection := client.Database("testing").Collection("hashKeyPairs")
	var result bson.M
	collection.FindOne(ctx, doc).Decode(&result)

	http.Redirect(w, r, result["url"].(string), 302)

	client.Disconnect(ctx)
}

func readFromDb(shortenedURL string, client *mongo.Client, ctx context.Context) string {
	doc := bson.M{
		"code": shortenedURL,
	}

	collection := client.Database("testing").Collection("hashKeyPairs")
	var result bson.M
	collection.FindOne(ctx, doc).Decode(&result)

	return result["url"].(string)
}

func writeToDb(content kvPair, client *mongo.Client, ctx context.Context) {
	doc := bson.M{
		"code": content.shortened,
		"url":  content.original,
		"hits": 0,
	}

	collection := client.Database("testing").Collection("hashKeyPairs")
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func getHash(s string) kvPair {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	sum := hash.Sum32()
	kv := kvPair{shortened: strconv.FormatUint(uint64(sum), 36), original: s}
	fmt.Println(kv.shortened)
	return kv
}
