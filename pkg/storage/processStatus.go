package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLastProcessStatus() (*models.ProcessStatus, error) {
	var processStatus models.ProcessStatus
	err := processStatusCollection.FindOne(context.Background(), bson.D{}).Decode(&processStatus)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &processStatus, nil
}

func UpdateProcessStatus(processStatus *models.ProcessStatus) error {
	filter := bson.M{"_id": processStatus.ID}
	update := bson.M{
		"$set": bson.M{
			"last_processed_file": processStatus.LastProcessedFile,
			"last_processed_time": processStatus.LastProcessedTime,
		},
	}

	_, err := processStatusCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func CreateProcessStatus(processStatus *models.ProcessStatus) error {
	result, err := processStatusCollection.InsertOne(context.Background(), processStatus)
	if err != nil {
		return err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("could not convert InsertedID to ObjectID")
	}
	processStatus.ID = oid

	return nil
}

// GetLastProcessedData returns the last processed file and date
func GetLastProcessedData() (*models.ProcessStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOne()
	var data models.ProcessStatus

	err := processStatusCollection.FindOne(ctx, bson.M{}, opts).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}
