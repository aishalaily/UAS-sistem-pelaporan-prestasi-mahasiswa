package repository

import (
	"context"
	"time"

	"uas-go/app/model"
	"uas-go/database"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func achievementCollection() *mongo.Collection {
	return database.MongoDB.Collection("achievements")
}

func InsertAchievementMongo(data model.AchievementMongo) (string, error) {
	id := uuid.NewString()
	data.ID = id
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	_, err := achievementCollection().InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetAchievementMongo(id string) (*model.AchievementMongo, error) {
	var res model.AchievementMongo
	err := achievementCollection().
		FindOne(context.Background(), bson.M{"_id": id}).
		Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func UpdateAchievementMongo(id string, update model.AchievementMongo) error {
	update.UpdatedAt = time.Now()
	_, err := achievementCollection().UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	return err
}
