package usecases

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/repositories"
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatBroadcaster defines the output port for the chat use case, allowing it to send
// messages without being aware of the underlying transport (e.g., WebSocket).
type ChatBroadcaster interface {
	BroadcastToRoom(roomID string, message []byte)
	SendMessage(clientID string, message []byte) error
}

// ChatUseCase defines the input port for chat-related business logic.
// It orchestrates operations like user connections, disconnections, and message processing.
type ChatUseCase interface {
	// UserConnected handles the logic when a new user joins a chat room.
	// It returns the recent chat history for the room.
	UserConnected(ctx context.Context, userID, roomID string) ([]*entities.MessageResponse, error)

	// UserDisconnected handles the logic when a user leaves a chat room.
	UserDisconnected(ctx context.Context, userID, roomID string) error

	// ProcessMessage handles an incoming message from a user, saves it, and broadcasts it.
	ProcessMessage(ctx context.Context, userID, roomID string, message []byte) error
}

type chatUseCase struct {
	userRepo    repositories.UserRepository
	messageRepo repositories.MessageRepository
	broadcaster ChatBroadcaster
}

func NewChatUseCase(
	userRepo repositories.UserRepository,
	messageRepo repositories.MessageRepository,
	broadcaster ChatBroadcaster,
) ChatUseCase {
	return &chatUseCase{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		broadcaster: broadcaster,
	}
}

// toMessageResponse converts a message entity to a message DTO, enriching it with user details.
func (uc *chatUseCase) toMessageResponse(ctx context.Context, msg *entities.Message) (*entities.MessageResponse, error) {
	var user *entities.User
	var err error

	// Handle system messages, which don't have a user in the repository.
	if msg.UserID == "system" {
		user = &entities.User{ID: "system", Username: "System"}
	} else {
		user, err = uc.userRepo.FindByID(ctx, msg.UserID)
		if err != nil {
			return nil, err
		}
	}

	return &entities.MessageResponse{
		ID:        msg.ID,
		Event:     "message",
		RoomID:    msg.RoomID,
		UserID:    msg.UserID,
		Username:  user.Username,
		UserRole:  user.Role,
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
		IsRead:    msg.IsRead,
		Type:      msg.Type,
		Metadata:  msg.Metadata,
	}, nil
}

// UserConnected handles new client connections.
func (uc *chatUseCase) UserConnected(ctx context.Context, userID, roomID string) ([]*entities.MessageResponse, error) {
	log.Printf("User %s connected to room %s", userID, roomID)

	messages, err := uc.messageRepo.FindByRoom(ctx, roomID)
	if err != nil {
		log.Printf("Failed to retrieve chat history for room %s: %v", roomID, err)
		return nil, nil
	}

	var history []*entities.MessageResponse
	for _, msg := range messages {
		dto, err := uc.toMessageResponse(ctx, msg)
		if err != nil {
			log.Printf("Failed to convert message to DTO: %v", err)
			continue
		}
		history = append(history, dto)
	}

	// Notify others that a user has joined.
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("Could not find user %s: %v", userID, err)
		return history, err
	}

	joinMsg := &entities.MessageResponse{
		ID:        primitive.NewObjectID(),
		Event:     "user-joined",
		RoomID:    roomID,
		UserID:    "system",
		Username:  "System",
		Content:   user.Username + " has joined the room.",
		Timestamp: time.Now(),
	}
	payload, _ := json.Marshal(joinMsg)
	uc.broadcaster.BroadcastToRoom(roomID, payload)

	return history, nil
}

// UserDisconnected handles client disconnections.
func (uc *chatUseCase) UserDisconnected(ctx context.Context, userID, roomID string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("Could not find user %s: %v", userID, err)
		return err
	}

	leaveMsg := &entities.MessageResponse{
		ID:        primitive.NewObjectID(),
		Event:     "user-left",
		RoomID:    roomID,
		UserID:    "system",
		Username:  "System",
		Content:   user.Username + " has left the room.",
		Timestamp: time.Now(),
	}
	payload, _ := json.Marshal(leaveMsg)
	uc.broadcaster.BroadcastToRoom(roomID, payload)

	return nil
}

// IncomingMessage represents the structure of a message received from a client.
type IncomingMessage struct {
	Type     string                 `json:"type"`
	Content  string                 `json:"content,omitempty"`
	Metadata *entities.FileMetadata `json:"metadata,omitempty"`
}

// ProcessMessage handles incoming chat messages.
func (uc *chatUseCase) ProcessMessage(ctx context.Context, userID, roomID string, rawMessage []byte) error {
	var incomingMsg IncomingMessage
	if err := json.Unmarshal(rawMessage, &incomingMsg); err != nil {
		log.Printf("Failed to unmarshal incoming message: %v", err)
		return err
	}

	// Create a new message and store it.
	msg := &entities.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    userID,
		Type:      incomingMsg.Type,
		Content:   incomingMsg.Content,  // For text messages
		Metadata:  incomingMsg.Metadata, // For file messages
		Timestamp: time.Now(),
		IsRead:    false,
	}

	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		log.Printf("Failed to save message to database: %v", err)
		return err
	}

	dto, err := uc.toMessageResponse(ctx, msg)
	if err != nil {
		log.Printf("Failed to create message DTO: %v", err)
		return err
	}

	payload, err := json.Marshal(dto)
	if err != nil {
		log.Printf("Failed to marshal message DTO: %v", err)
		return err
	}

	uc.broadcaster.BroadcastToRoom(roomID, payload)
	return nil
}
