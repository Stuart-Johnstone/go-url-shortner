package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type kvPair struct {
	shortened string
	original  string
}

var tmpl = template.Must(template.ParseFiles("templates/home.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", handleRedirect)
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /shorten", handleShorten)
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

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
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

func handleShorten(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := getHash(r.FormValue("url"))
	doc := bson.M{
		"code": hash.shortened,
		"url":  r.FormValue("url"),
		"hits": 0,
	}
	client, ctx := dbConnect()
	collection := client.Database("testing").Collection("hashKeyPairs")
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, map[string]string{
		"ShortURL": ("http://localhost:8080/" + hash.shortened),
	})
	client.Disconnect(ctx)
	fmt.Println(result)
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

	url := result["url"].(string)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	http.Redirect(w, r, url, 302)

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
