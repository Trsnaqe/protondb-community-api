package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty"`
	AppID   string               `bson:"appId"`
	Title   *string              `bson:"title"`
	Reports []primitive.ObjectID `bson:"reports"`
}

func NewGame(appID string, title *string) *Game {
	return &Game{
		AppID:   appID,
		Title:   title,
		Reports: []primitive.ObjectID{},
	}
}
