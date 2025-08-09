package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ErrRecordNotFound = errors.New("record not found")

type taskRepository struct {
	collection *mongo.Collection
	db         *mongo.Client
}

func NewTaskRepository(conn dbContract.DBConnection) repoContract.TaskRepository {
	coll := conn.DB().Database("microgate").Collection("task")
	return &taskRepository{
		collection: coll,
		db:         conn.DB(),
	}
}

func (t *taskRepository) CreateTask(ctx context.Context, task *entity.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	e := bson.D{
		{"user_id", task.UserID},
		{"title", task.Title},
		{"description", task.Description},
		{"status", task.Status},
	}

	i, err := t.collection.InsertOne(ctx, e)
	if err != nil {
		return fmt.Errorf("failed to insert task %s", err.Error())
	}
	task.Id = i.InsertedID.(primitive.ObjectID)
	return nil
}

func (t *taskRepository) ListTask(ctx context.Context, userID int64) ([]entity.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cursor, err := t.collection.Find(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks %s", err.Error())
	}
	if cursor.Err() != nil {
		return nil, fmt.Errorf("failed to get tasks from db %w", err)
	}
	var tasks []entity.Task
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.New("failed to get events from database")
	}

	for _, res := range results {
		task := entity.Task{
			Id:          res["_id"].(primitive.ObjectID),
			UserID:      res["user_id"].(int64),
			Title:       res["title"].(string),
			Description: res["description"].(string),
			Status:      res["status"].(string),
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

func (t *taskRepository) DeleteTask(ctx context.Context, taskId string, userID int64) error {
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return errors.New("failed to convert eventId to ObjectID")
	}
	return t.collection.FindOneAndDelete(ctx, bson.M{"user_id": userID, "_id": id}).Err()
}

func (t *taskRepository) ServiceStatus(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return t.db.Ping(ctx, readpref.Primary())
}
