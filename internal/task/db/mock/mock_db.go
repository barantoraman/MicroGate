package mock

import (
	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockConnection struct{}

func NewMockConnection() (dbContract.DBConnection, error) {
	return &mockConnection{}, nil
}

func (m *mockConnection) DB() *mongo.Client { return nil }

func (m *mockConnection) Close() {}
