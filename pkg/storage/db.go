package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client                  *mongo.Client
	gamesCollection         *mongo.Collection
	reportsCollection       *mongo.Collection
	processStatusCollection *mongo.Collection
)

func ConnectDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	dbURI := os.Getenv("DB_URI")
	opts := options.Client().ApplyURI(dbURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database("reports").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	gamesCollection = client.Database("protondb_reports").Collection("games")
	reportsCollection = client.Database("protondb_reports").Collection("reports")
	processStatusCollection = client.Database("protondb_reports").Collection("process_status")

	// Ensure the index on the title field
	if err := ensureTitleIndex(); err != nil {
		log.Printf("Error creating index: %v", err)
		return err
	}

	return nil
}

func CloseDB() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}
}

// ensureTitleIndex creates an index on the title field of the games collection if it doesn't exist
func ensureTitleIndex() error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: "text"}},
		Options: options.Index().SetUnique(false),
	}

	_, err := gamesCollection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}
