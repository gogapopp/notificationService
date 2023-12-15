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
			Port:     "21017",
		},
	}
	_, err := NewMongoDB(cfg)
	assert.Error(t, err, "failed to create DB")
}
