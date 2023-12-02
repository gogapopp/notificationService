package models

import "time"

type Message struct {
	UserID    string    `json:"user_id" bson:"user_id" validate:"required"`
	Message   string    `json:"message" bson:"message" validate:"required"`
	Timestamp time.Time `json:"-"`
}

type UserSub struct {
	UserID string `json:"user_id" bson:"user_id" validate:"required"`
	Email  string `json:"email" bson:"email" validate:"required,email"`
}
