package mongodb

import (
	"context"
	"fmt"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/gogapopp/notificationService/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbName        = "mongodb"
	collectionMessages = "messages"
	collectionUsers    = "users"
)

type DB struct {
	Client *mongo.Client
}

func NewMongoDB(config *config.Config) (*DB, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.MongoDB.User, config.MongoDB.Password, config.MongoDB.Host, config.MongoDB.Port)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &DB{Client: client}, client.Ping(ctx, nil)
}

func (d *DB) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx, nil)
}

func (d *DB) InsertMessage(ctx context.Context, msg models.Message) error {
	collection := d.Client.Database(mongodbName).Collection(collectionMessages)
	_, err := collection.InsertOne(ctx, bson.M{
		"user_id":   msg.UserID,
		"message":   msg.Message,
		"timestamp": msg.Timestamp,
	})
	return err
}

func (d *DB) Subscribe(ctx context.Context, userSub models.UserSub) error {
	collection := d.Client.Database(mongodbName).Collection(collectionUsers)
	_, err := collection.InsertOne(ctx, bson.M{
		"user_id": userSub.UserID,
		"email":   userSub.Email,
	})
	return err
}
