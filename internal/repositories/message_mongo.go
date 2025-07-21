package repositories

import (
	"api-gateway/internal/entities"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageRepository defines the interface for message data storage.
type MessageRepository interface {
	// Create stores a new message in the database.
	Create(ctx context.Context, message *entities.Message) error
	// FindByRoom retrieves all messages for a given room, sorted by timestamp.
	FindByRoom(ctx context.Context, roomID string) ([]*entities.Message, error)
}

// mongoMessageRepository is a MongoDB implementation of the MessageRepository.
type mongoMessageRepository struct {
	collection *mongo.Collection
}

// NewMongoMessageRepository creates a new MongoDB message repository.
func NewMongoMessageRepository(db *mongo.Database) MessageRepository {
	return &mongoMessageRepository{
		collection: db.Collection("messages"),
	}
}

// Create inserts a new message into the MongoDB collection.
func (r *mongoMessageRepository) Create(ctx context.Context, message *entities.Message) error {
	log.Println("create message")
	_, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		log.Println("create message error: ", err)
		return err
	}

	return nil
}

// FindByRoom retrieves all messages for a given room, sorted by timestamp.
func (r *mongoMessageRepository) FindByRoom(ctx context.Context, roomID string) ([]*entities.Message, error) {
	filter := bson.M{"room_id": roomID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
