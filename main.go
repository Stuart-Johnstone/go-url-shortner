package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"

	"urlshortener/util"
)

var tmpl = template.Must(template.ParseFiles("templates/home.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", handleRedirect)
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /shorten", handleShorten)

	http.ListenAndServe(":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	hash := util.GetHash(r.FormValue("url"))

	client, ctx := util.DbConnect()

	collection := client.Database("testing").Collection("hashKeyPairs")

	doc := bson.M{"code": hash.shortened, "url": r.FormValue("url"), "hits": 0}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, map[string]string{
		"ShortURL": ("http://localhost:8080/" + hash.shortened),
	})
	client.Disconnect(ctx)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	client, ctx := util.DbConnect()

	path := r.PathValue("id")
	doc := bson.M{"code": path}

	var result bson.M
	collection := client.Database("testing").Collection("hashKeyPairs")
	collection.FindOne(ctx, doc).Decode(&result)

	url := result["url"].(string)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	http.Redirect(w, r, url, 302)

	update := bson.M{"$inc": bson.M{"hits": 1}}
	collection.UpdateOne(ctx, doc, update)

	client.Disconnect(ctx)
}
