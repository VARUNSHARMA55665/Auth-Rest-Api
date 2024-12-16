package db

import (
	"auth-rest-api/constants"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitMongoClient() error {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mongoBase := constants.MongoBase
	mongoPass := constants.MongoPass
	mongoUri := constants.MongoUri

	mongoBaseURI := mongoBase + ":" + mongoPass + "@" + mongoUri

	fmt.Println(mongoBaseURI)

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoBaseURI))

	if err == nil {
		log.Println("InitMongoClient DB Connect Success")
	}

	return err
}

func GetMongoCollection(collectionName string) *mongo.Collection {
	dbName := constants.DbName
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

func FindOneMongo(collectionName string, filter interface{}, result interface{}) error {
	collection := GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	return err
}

func UpdateOneMongo(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	collection := GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update, opts...)
	return err
}

func FindAllMongo(collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	collection := GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, opts...)
	return cursor, err
}
