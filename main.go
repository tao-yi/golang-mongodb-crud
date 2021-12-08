package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Person struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func CreatePersonEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var person Person
	err := json.NewDecoder(req.Body).Decode(&person)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}

	collection := client.Database("test").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	json.NewEncoder(res).Encode(result)
}

func GetPeopleEndpoint(res http.ResponseWriter, request *http.Request) {
	res.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("test").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	json.NewEncoder(res).Encode(people)
}

func GetPersonEndpoint(res http.ResponseWriter, request *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("test").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var person Person
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	json.NewEncoder(res).Encode(person)
}

func main() {
	fmt.Println("Starting the application")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootpassword@localhost:27017"))
	// Make sure to defer a call to Disconnect after instantiating your client:
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods(http.MethodPost)
	router.HandleFunc("/person", GetPeopleEndpoint).Methods(http.MethodGet)
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
