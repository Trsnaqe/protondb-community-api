package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/trsnaqe/protondb-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateReport(report *models.Report, game *models.Game) (*models.Report, error) {
	result, err := reportsCollection.InsertOne(context.Background(), report)
	if err != nil {
		return nil, err
	}

	report.ID = result.InsertedID.(primitive.ObjectID)

	return report, nil
}

func GetReportByID(reportID string) (*models.Report, error) {
	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var report models.Report
	err = reportsCollection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &report, nil
}

func GetAllReports() (*mongo.Cursor, error) {
	cursor, err := reportsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func UpdateReport(report *models.Report) error {
	if report == nil {
		return errors.New("nil report provided")
	}

	filter := bson.M{"_id": report.ID}
	update := bson.M{"$set": bson.M{"data": report.Data, "report_version": report.ReportVersion}}

	_, err := reportsCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func GetReportsByGameID(gameID string, version string) ([]models.Report, error) {
	if gameID == "" {
		log.Println("Error: empty gameID provided")
		return nil, errors.New("empty gameID provided")
	}

	game, err := GetGameByAppIDWithReports(gameID)
	if err != nil {
		log.Println("Error getting game by AppID:", err)
		return nil, err
	}

	if game == nil || game.Reports == nil || len(game.Reports) == 0 {
		log.Println("No reports found for gameID:", gameID)
		return []models.Report{}, nil
	}

	reportsFilter := bson.M{"_id": bson.M{"$in": game.Reports}}

	if version == "V1" || version == "V2" {
		reportsFilter["report_version"] = version
	}

	var reports []models.Report
	cursor, err := reportsCollection.Find(context.Background(), reportsFilter)
	if err != nil {
		log.Println("Error finding reports:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Println("Error on closing cursor:", err)
		}
	}()

	for cursor.Next(context.Background()) {
		var report models.Report
		if err := cursor.Decode(&report); err != nil {
			log.Println("Error decoding report:", err)
			return nil, fmt.Errorf("error decoding report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error reading from cursor:", err)
		return nil, fmt.Errorf("error reading from cursor: %w", err)
	}

	return reports, nil
}

func DeleteReport(reportID string) error {
	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = reportsCollection.DeleteOne(context.Background(), filter)
	return err
}

// GetTotalReportsCount returns the total number of reports
func GetTotalReportsCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := reportsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return count, nil
}
func CompareReport(report models.ReportFormatV2) bool {

	filter := bson.D{
		{Key: "data.app.steam.appId", Value: report.App.Steam.AppID},
		{Key: "data.app.title", Value: report.App.Title},
		{Key: "data.timestamp", Value: report.Timestamp},
		{Key: "data.systemInfo.cpu", Value: report.SystemInfo.CPU},
		{Key: "data.systemInfo.gpu", Value: report.SystemInfo.GPU},
		{Key: "data.systemInfo.gpuDriver", Value: report.SystemInfo.GPUDriver},
		{Key: "data.systemInfo.kernel", Value: report.SystemInfo.Kernel},
		{Key: "data.systemInfo.os", Value: report.SystemInfo.OS},
		{Key: "data.systemInfo.ram", Value: report.SystemInfo.RAM},
	}

	var result models.ReportFormatV2
	err := reportsCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Printf("Could not perform find operation: %v", err)
		return false
	}
	return true
}
func CountV2Reports() (int64, error) {
	filter := bson.D{{Key: "report_version", Value: "V2"}}
	count, err := reportsCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Printf("Error while counting V2 reports: %v", err)
		return 0, err
	}
	return count, nil
}

// search game by title and get its reports
func GetReportsOfMatchedGamesByTitle(title string, versioned bool, version string) (*mongo.Cursor, error) {
	filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}
	cursor, err := gamesCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var matchedGames []models.Game
	for cursor.Next(context.Background()) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			log.Printf("Error decoding game: %v", err)
			continue
		}
		matchedGames = append(matchedGames, game)
	}

	cursor.Close(context.Background())

	var gameIDs []primitive.ObjectID
	for _, game := range matchedGames {
		gameIDs = append(gameIDs, game.ID)
	}

	reportsFilter := bson.M{
		"_id": bson.M{"$in": gameIDs},
	}

	if version == "V1" || version == "V2" {
		reportsFilter["report_version"] = version
	}

	options := options.Find().SetProjection(bson.M{"reports": 1}) // Adjust projection as needed

	reportsCursor, err := reportsCollection.Find(context.Background(), reportsFilter, options)
	if err != nil {
		return nil, err
	}

	return reportsCursor, nil
}
