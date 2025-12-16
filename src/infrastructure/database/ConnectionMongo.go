package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionMongo() (*mongo.Client, error) {

	err := godotenv.Load(".env.dev")
	if err != nil {
		panic("Error loading .env file")
	}
	var uri string
	if uri = os.Getenv("MONGODB_URI_ADMIN"); uri == "" {
		panic("MONGODB_URI not set in environment")
	}
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	err = client.Database("luminaMO-orq-modAI").CreateCollection(context.TODO(), "testCollection")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil

}
