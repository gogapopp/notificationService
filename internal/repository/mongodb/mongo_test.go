package mongodb

import (
	"testing"

	"github.com/gogapopp/notificationService/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMongoDB(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDB{
			User:     "testUser",
			Password: "testPass",
			Host:     "host",
			Port:     "27017",
		},
	}
	db, err := NewMongoDB(cfg)
	assert.NoError(t, err, "failed to create DB")
	assert.NotNil(t, db.Client, "expected DB client not nil")
}
