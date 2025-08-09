package mongo

import (
	"context"
	"time"

	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	"github.com/barantoraman/microgate/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoConnection struct {
	db *mongo.Client
}

func NewMongoConnection(cfg config.TaskServiceConfigurations) (dbContract.DBConnection, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &mongoConnection{
		db: client,
	}, nil
}

func (m *mongoConnection) Close() {
	m.db.Disconnect(context.TODO())
}

func (m *mongoConnection) DB() *mongo.Client {
	return m.db
}
