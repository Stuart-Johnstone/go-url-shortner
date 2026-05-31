package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"urlshortener/model"
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

	doc := bson.M{"code": hash.Shortened, "url": r.FormValue("url"), "hits": 0}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, map[string]string{
		"ShortURL": ("http://localhost:8080/" + hash.Shortened),
	})
	client.Disconnect(ctx)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	client, ctx := util.DbConnect()
	collection := client.Database("testing").Collection("hashKeyPairs")

	path := r.PathValue("id")
	doc := bson.M{"code": path}

	res, err := util.CheckCache(path)
	fmt.Println("Cache Result: " + res)
	if err != nil {

		var result bson.M
		collection.FindOne(ctx, doc).Decode(&result)
		res = result["url"].(string)

		fmt.Println("Cache missed, writing " + res + " to cache")
		util.WriteCache(model.KvPair{Shortened: path, Original: res})
	}

	if !strings.HasPrefix(res, "http://") && !strings.HasPrefix(res, "https://") {
		res = "https://" + res
	}

	http.Redirect(w, r, res, 302)

	update := bson.M{"$inc": bson.M{"hits": 1}}
	collection.UpdateOne(ctx, doc, update)

	client.Disconnect(ctx)
}
