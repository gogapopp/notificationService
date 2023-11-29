package mongodb

import (
	"context"
	"fmt"

	"github.com/gogapopp/notificationService/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
