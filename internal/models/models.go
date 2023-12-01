package models

type Message struct {
	UserID  string `json:"user_id" bson:"user_id"`
	Message string `json:"message" bson:"message"`
}
