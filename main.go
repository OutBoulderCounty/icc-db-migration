package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoOption struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Index int                `bson:"index"`
}

type mongoElement struct {
	ID         primitive.ObjectID   `bson:"_id"`
	Label      string               `bson:"label"`
	Type       string               `bson:"type"`
	Index      int                  `bson:"index"`
	Required   bool                 `bson:"required"`
	Categories []primitive.ObjectID `bson:"categories"`
	Options    []mongoOption        `bson:"options"`
	Priority   int                  `bson:"priority"`
	Search     bool                 `bson:"search"`
}

type mongoForm struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Elements []mongoElement     `bson:"elements"`
	Required bool               `bson:"required"`
	Live     bool               `bson:"live"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env: " + err.Error())
	}
	// get forms data from prod mongodb deployment
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: " + err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: " + err.Error())
	}
	fmt.Println("MongoDB connection successful")
	defer client.Disconnect(ctx)
	db := client.Database("healthdir")
	formsColl := db.Collection("forms")
	filter := bson.M{}
	cursor, err := formsColl.Find(ctx, filter)
	if err != nil {
		log.Fatal("Failed to retrieve forms: " + err.Error())
	}
	var forms []*mongoForm
	err = cursor.All(ctx, &forms)
	if err != nil {
		log.Fatal("Failed to parse forms: " + err.Error())
	}
	for i := 0; i < len(forms); i++ {
		form := forms[i]
		fmt.Println(*form)
	}

	// insert forms data into singlestore cluster

	// get all data from local mongodb

	// insert data into singlestore
}
