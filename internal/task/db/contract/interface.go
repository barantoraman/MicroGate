package contract

import "go.mongodb.org/mongo-driver/mongo"

type DBConnection interface {
	DB() *mongo.Client
	Close()
}
