package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/gogapopp/notificationService/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongodbName = "mongodb"

type DB struct {
	Client *mongo.Client
}

func NewMongoDB(config *config.Config) (*DB, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.MongoDB.User, config.MongoDB.Password, config.MongoDB.Host, config.MongoDB.Port)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &DB{Client: client}, nil
}

func (d *DB) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx, nil)
}

func (d *DB) InsertMessage(ctx context.Context, msg models.Message) error {
	collection := d.Client.Database(mongodbName).Collection("messages")
	_, err := collection.InsertOne(ctx, bson.M{
		"user_id":   msg.UserID,
		"message":   msg.Message,
		"timestamp": time.Now(),
	})
	return err
}
