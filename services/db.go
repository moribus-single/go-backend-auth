package services

import (
	"app/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbService struct {
	User *mongo.Collection
}

func GetDbService(db, table, uri string) (DbService, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.TODO()

	// Connect to the mongo database
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return DbService{}, nil
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return DbService{}, err
	}

	return DbService{
		User: client.Database(db).Collection(table),
	}, nil
}

func (s DbService) Read(guid string) (models.User, error) {
	id, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return models.User{}, err
	}

	var result models.User
	filter := bson.D{{"_id", id}}

	err = s.User.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.User{}, err
	}

	return result, nil
}

func (s DbService) Update(guid, refresh string) error {
	instance, err := s.Read(guid)
	if err != nil {
		log.Fatal(err)
	}

	instance.Refresh = refresh
	instance.TokenCounter += 1

	id, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{"_id", id}}
	s.User.FindOneAndReplace(context.TODO(), filter, instance)

	return nil
}
