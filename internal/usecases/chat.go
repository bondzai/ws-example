package usecases

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/repositories"
	"api-gateway/pkg/ws"
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatUseCase defines the business logic for the chat system.
type ChatUseCase interface {
	// OnConnect is called when a new client connects.
	OnConnect(c *ws.Client)
	// OnDisconnect is called when a client disconnects.
	OnDisconnect(c *ws.Client)
	// OnMessage is called when a message is received from a client.
	OnMessage(c *ws.Client, messageType int, data []byte)
}

// chatUseCase is the implementation of the ChatUseCase.
type chatUseCase struct {
	userRepo    repositories.UserRepository
	messageRepo repositories.MessageRepository
	wsHandler   *ws.WSHandler
}

// NewChatUseCase creates a new ChatUseCase.
func NewChatUseCase(userRepo repositories.UserRepository, messageRepo repositories.MessageRepository, wsHandler *ws.WSHandler) ChatUseCase {
	uc := &chatUseCase{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		wsHandler:   wsHandler,
	}

	// Register the use case methods as the WebSocket handler callbacks.
	wsHandler.OnConnect = uc.OnConnect
	wsHandler.OnMessage = uc.OnMessage
	wsHandler.OnClose = uc.OnDisconnect

	return uc
}

// OnConnect handles new client connections.
func (uc *chatUseCase) OnConnect(c *ws.Client) {
	log.Printf("Client %s connected to room %s", c.GetID(), c.GetRoomID())

	// Retrieve chat history for the room.
	messages, err := uc.messageRepo.FindByRoom(context.Background(), c.GetRoomID())
	if err != nil {
		log.Printf("Failed to retrieve chat history for room %s: %v", c.GetRoomID(), err)
		// We don't return here, as we can still proceed with the connection.
	}

	// Send the history to the newly connected client.
	for _, msg := range messages {
		payload, _ := json.Marshal(msg)
		c.SendMessage(payload)
	}

	// Notify others in the room that a new user has joined.
	user, err := uc.userRepo.FindByID(context.Background(), c.GetID())
	if err != nil {
		log.Printf("Could not find user %s: %v", c.GetID(), err)
		return
	}

	joinMsg := entities.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    c.GetRoomID(),
		UserID:    user.ID,
		Username:  "System",
		Content:   user.Username + " has joined the room.",
		Timestamp: time.Now(),
	}
	payload, _ := json.Marshal(joinMsg)

	// If auto-sync is enabled, use the broker to broadcast. Otherwise, broadcast locally.
	if uc.wsHandler.GetConfig().EnableAutoSync {
		uc.wsHandler.SyncMessage(context.Background(), c, payload)
	} else {
		uc.wsHandler.BroadcastToRoom(c.GetRoomID(), payload)
	}
}

// OnDisconnect handles client disconnections.
func (uc *chatUseCase) OnDisconnect(c *ws.Client) {
	log.Printf("Client %s disconnected from room %s", c.GetID(), c.GetRoomID())

	// Notify others in the room that a user has left.
	user, err := uc.userRepo.FindByID(context.Background(), c.GetID())
	if err != nil {
		log.Printf("Could not find user %s: %v", c.GetID(), err)
		return
	}

	leaveMsg := entities.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    c.GetRoomID(),
		UserID:    user.ID,
		Username:  "System",
		Content:   user.Username + " has left the room.",
		Timestamp: time.Now(),
	}
	payload, _ := json.Marshal(leaveMsg)

	// If auto-sync is enabled, use the broker. Otherwise, broadcast locally.
	if uc.wsHandler.GetConfig().EnableAutoSync {
		uc.wsHandler.SyncMessage(context.Background(), c, payload)
	} else {
		uc.wsHandler.BroadcastToRoom(c.GetRoomID(), payload)
	}
}

// OnMessage handles incoming chat messages.
func (uc *chatUseCase) OnMessage(c *ws.Client, messageType int, data []byte) {
	user, err := uc.userRepo.FindByID(context.Background(), c.GetID())
	if err != nil {
		log.Printf("Could not find user %s: %v", c.GetID(), err)
		return
	}

	// Create a new message and store it.
	msg := &entities.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    c.GetRoomID(),
		UserID:    c.GetID(),
		Username:  user.Username,
		UserRole:  user.Role,
		Content:   string(data),
		Timestamp: time.Now(),
		IsRead:    false, // Default to false
	}

	if err := uc.messageRepo.Create(context.Background(), msg); err != nil {
		log.Printf("Failed to save message to database: %v", err)
		return
	}

	// Marshal the message to JSON to be sent to clients.
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// If auto-sync is on, publish to Redis. Otherwise, broadcast locally.
	if uc.wsHandler.GetConfig().EnableAutoSync {
		uc.wsHandler.SyncMessage(context.Background(), c, payload)
	} else {
		uc.wsHandler.BroadcastToRoom(c.GetRoomID(), payload)
	}
}
