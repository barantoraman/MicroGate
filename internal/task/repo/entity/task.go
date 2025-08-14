package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id          primitive.ObjectID `json:"id"`
	UserID      int64              `json:"user_id"`
	Title       string             `json:"title"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status"`
}
