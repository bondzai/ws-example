package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a single chat message.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"room_id" json:"roomId"`
	UserID    string             `bson:"user_id" json:"userId"`
	Username  string             `bson:"username" json:"username"`
	UserRole  UserRole           `bson:"user_role" json:"userRole"`
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	IsRead    bool               `bson:"is_read" json:"isRead"`
}
