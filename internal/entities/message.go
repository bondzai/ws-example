package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a single chat message as stored in the database.
// It only contains the UserID to avoid data duplication.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"room_id" json:"roomId"`
	UserID    string             `bson:"user_id" json:"userId"`
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	IsRead    bool               `bson:"is_read" json:"isRead"`
	Type      string             `bson:"type" json:"type"` // "text" or "file"
	Metadata  *FileMetadata      `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// FileMetadata holds information about an uploaded file.
type FileMetadata struct {
	URL      string `bson:"url" json:"url"`
	FileName string `bson:"file_name" json:"fileName"`
	FileSize int64  `bson:"file_size" json:"fileSize"`
	MIMEType string `bson:"mime_type" json:"mimeType"`
}

// MessageResponse is a DTO for sending message data to clients.
// It includes user details, which are populated at runtime.
type MessageResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Event     string             `json:"event"`
	RoomID    string             `json:"roomId"`
	UserID    string             `json:"userId"`
	Username  string             `json:"username"`
	UserRole  UserRole           `json:"userRole"`
	Content   string             `json:"content"`
	Timestamp time.Time          `json:"timestamp"`
	IsRead    bool               `json:"isRead"`
	Type      string             `json:"type"`
	Metadata  *FileMetadata      `json:"metadata,omitempty"`
}
