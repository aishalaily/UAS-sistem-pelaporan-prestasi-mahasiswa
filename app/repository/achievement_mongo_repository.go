package repository

import (
	"context"
	"time"
	"uas-go/app/model"
	"uas-go/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertAchievementMongo(studentID string, req model.AchievementRequest) (string, error) {
	collection := database.MongoDB.Collection("achievements")

	doc := bson.M{
		"student_id": studentID,
		"title": req.Title,
		"category": req.Category,
		"description": req.Description,
		"event_date": req.EventDate,
		"documents": req.Documents,
		"status": "draft",
		"created_at": time.Now(),
	}

	res, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func UpdateAchievementMongo(col *mongo.Collection, mongoID string, req model.AchievementRequest) error {
    update := bson.M{
        "$set": bson.M{
            "title": req.Title,
			"category": req.Category,
			"description": req.Description,
			"event_date": req.EventDate,
			"documents": req.Documents,
            "updated_at":  time.Now(),
        },
    }

    _, err := col.UpdateByID(context.Background(), mongoID, update)
    return err
}


func SoftDeleteAchievementMongo(mongoID string) error {
	collection := database.MongoDB.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		update,
	)

	return err
}



