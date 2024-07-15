package storage

import (
	"context"
	"log"
	"time"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateGame(game *models.Game) error {
	if game.Reports == nil {
		game.Reports = []primitive.ObjectID{}
	}
	result, err := gamesCollection.InsertOne(context.Background(), game)
	if err != nil {
		return err
	}
	game.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func GetGameByAppID(appID string) (*models.Game, error) {
	filter := bson.M{"appId": appID}
	options := options.FindOne().SetProjection(bson.M{"reports": 0})

	var game models.Game
	err := gamesCollection.FindOne(context.Background(), filter, options).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &game, nil
}

func GetGameByAppIDWithReports(appID string) (*models.Game, error) {
	filter := bson.M{"appId": appID}

	var game models.Game
	err := gamesCollection.FindOne(context.Background(), filter).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &game, nil
}

func GetTotalGamesCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := gamesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return count, nil
}

func UpdateGame(game *models.Game) error {
	filter := bson.M{"_id": game.ID}
	update := bson.M{"$set": bson.M{"appId": game.AppID, "title": game.Title, "reports": game.Reports}}

	_, err := gamesCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func GetGameByID(gameID string) (*models.Game, error) {
	objectID, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	options := options.FindOne().SetProjection(bson.M{"reports": 0})

	var game models.Game
	err = gamesCollection.FindOne(context.Background(), filter, options).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &game, nil
}

func GetAllGames() (*mongo.Cursor, error) {
	options := options.Find().SetProjection(bson.M{"reports": 0})

	cursor, err := gamesCollection.Find(context.Background(), bson.M{}, options)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func DeleteGame(gameID string) error {
	objectID, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = gamesCollection.DeleteOne(context.Background(), filter)
	return err
}
func InsertReport(game *models.Game, reportID primitive.ObjectID) error {
	filter := bson.M{"_id": game.ID}

	checkUpdate := bson.M{"$setOnInsert": bson.M{"reports": bson.A{}}}
	_, err := gamesCollection.UpdateOne(context.Background(), filter, checkUpdate, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	update := bson.M{"$push": bson.M{"reports": reportID}}
	_, err = gamesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func ChangeTitle(gameID primitive.ObjectID, newTitle string) error {
	filter := bson.M{"_id": gameID}
	update := bson.M{"$set": bson.M{"title": newTitle}}

	_, err := gamesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// search game bby titles
func SearchGameByTitle(title string) (*mongo.Cursor, error) {
	filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}
	cursor, err := gamesCollection.Find(context.Background(), filter)
	return cursor, err
}

func GetGameByTitle(title string) (*models.Game, error) {
	filter := bson.M{"title": title}
	var game models.Game
	err := gamesCollection.FindOne(context.Background(), filter).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}
