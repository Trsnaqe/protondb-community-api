// located in pkg/models/process_status.go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessStatus struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	LastProcessedFile string             `bson:"last_processed_file"`
	LastProcessedTime primitive.DateTime `bson:"last_processed_time"`
}
