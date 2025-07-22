package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"
)

// Room holds a set of clients in a chat room.
type Room struct {
	Clients map[*Client]bool
}

// Hub maintains the set of active rooms and broadcasts messages to rooms.
type Hub struct {
	Rooms         map[string]*Room
	ClientsByID   map[string]*Client // Added for direct messaging
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan RoomMessage
	DirectMessage chan DirectMessage // Added for sending to a specific client
	mu            sync.Mutex
}

// RoomMessage is a message to be broadcast to a specific room.
type RoomMessage struct {
	RoomID  string
	Message []byte
}

// DirectMessage is a message to be sent to a specific client.
type DirectMessage struct {
	ClientID string
	Message  []byte
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Rooms:         make(map[string]*Room),
		ClientsByID:   make(map[string]*Client),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan RoomMessage),
		DirectMessage: make(chan DirectMessage),
	}
}

// Run starts the hub loop to process client registration, unregistration, and broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			// Register client to a room
			room, ok := h.Rooms[client.RoomID]
			if !ok {
				room = &Room{Clients: make(map[*Client]bool)}
				h.Rooms[client.RoomID] = room
			}
			room.Clients[client] = true
			// Register client by ID for direct messaging
			h.ClientsByID[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if room, ok := h.Rooms[client.RoomID]; ok {
				if _, exists := room.Clients[client]; exists {
					delete(room.Clients, client)
					close(client.Send)
					if len(room.Clients) == 0 {
						delete(h.Rooms, client.RoomID)
					}
				}
			}
			// Unregister client by ID
			delete(h.ClientsByID, client.ID)
			h.mu.Unlock()

		case roomMsg := <-h.Broadcast:
			h.mu.Lock()
			if room, ok := h.Rooms[roomMsg.RoomID]; ok {
				for client := range room.Clients {
					client.SendMessage(roomMsg.Message)
				}
			}
			h.mu.Unlock()

		case directMsg := <-h.DirectMessage:
			h.mu.Lock()
			if client, ok := h.ClientsByID[directMsg.ClientID]; ok {
				client.SendMessage(directMsg.Message)
			}
			h.mu.Unlock()
		}
	}
}

// WSConfig holds WebSocket configuration
type WSConfig struct {
	PingInterval   time.Duration
	PongWait       time.Duration
	WriteWait      time.Duration
	MaxMessageSize int64
	BufferSize     int
	EnableAutoSync bool
	SyncChannel    string
}

// WSHandler provides a hub for WebSocket connections and acts as a broadcaster.
type WSHandler struct {
	hub    *Hub
	config WSConfig
	broker MessageBroker
}

// SyncMessage defines the structure for synchronization messages.
type SyncMessage struct {
	ClientID string `json:"client_id,omitempty"` // Can be empty if it's a room broadcast
	RoomID   string `json:"room_id"`
	Data     []byte `json:"data"`
}

// NewWSHandler initializes a new WSHandler with its hub and message broker.
func NewWSHandler(opts ...Option) *WSHandler {
	// Start with a default configuration
	handler := &WSHandler{
		hub:    NewHub(),
		broker: NewNoOpMessageBroker(), // Default to no-op broker
		config: WSConfig{
			PingInterval:   30 * time.Second,
			PongWait:       60 * time.Second,
			WriteWait:      10 * time.Second,
			MaxMessageSize: 512,
			BufferSize:     256,
			EnableAutoSync: false,
			SyncChannel:    "websocket_sync",
		},
	}

	// Apply all the functional options to customize the handler
	for _, opt := range opts {
		opt(handler)
	}

	go handler.hub.Run()

	// Subscribe to sync messages if auto-sync is enabled.
	if handler.config.EnableAutoSync {
		if handler.broker.GetType() == "noop" {
			log.Println("Warning: Auto-sync is enabled, but no message broker is configured. Sync will not work.")
		} else {
			err := handler.broker.Subscribe(context.Background(), handler.config.SyncChannel, handler.handleSyncMessage)
			if err != nil {
				log.Printf("Failed to subscribe to sync channel: %v", err)
			}
		}
	}

	return handler
}

// handleSyncMessage processes incoming sync messages from other nodes
func (h *WSHandler) handleSyncMessage(data []byte) {
	var syncMsg SyncMessage
	if err := json.Unmarshal(data, &syncMsg); err != nil {
		log.Printf("Failed to unmarshal sync message: %v", err)
		return
	}

	// If ClientID is present, it's a direct message. Otherwise, broadcast to the room.
	if syncMsg.ClientID != "" {
		h.hub.DirectMessage <- DirectMessage{ClientID: syncMsg.ClientID, Message: syncMsg.Data}
	} else {
		h.hub.Broadcast <- RoomMessage{RoomID: syncMsg.RoomID, Message: syncMsg.Data}
	}
}

// RegisterClient registers a client with the hub.
func (h *WSHandler) RegisterClient(client *Client) {
	h.hub.Register <- client
}

// UnregisterClient unregisters a client from the hub.
func (h *WSHandler) UnregisterClient(client *Client) {
	h.hub.Unregister <- client
}

// BroadcastToRoom sends a message to all clients in a specific room.
// If auto-sync is enabled, it publishes the message to the message broker.
func (h *WSHandler) BroadcastToRoom(roomID string, message []byte) {
	if h.config.EnableAutoSync {
		syncMsg := SyncMessage{RoomID: roomID, Data: message}
		if err := h.broker.Publish(context.Background(), h.config.SyncChannel, syncMsg); err != nil {
			log.Printf("Failed to publish sync message: %v", err)
		}
	} else {
		h.hub.Broadcast <- RoomMessage{RoomID: roomID, Message: message}
	}
}

// SendMessage sends a message directly to a specific client by their ID.
func (h *WSHandler) SendMessage(clientID string, message []byte) error {
	if h.config.EnableAutoSync {
		syncMsg := SyncMessage{ClientID: clientID, RoomID: "", Data: message} // RoomID can be empty
		if err := h.broker.Publish(context.Background(), h.config.SyncChannel, syncMsg); err != nil {
			log.Printf("Failed to publish direct sync message: %v", err)
			return err
		}
	} else {
		// Check if client is local before sending
		h.hub.mu.Lock()
		_, ok := h.hub.ClientsByID[clientID]
		h.hub.mu.Unlock()
		if !ok {
			return errors.New("client not found")
		}
		h.hub.DirectMessage <- DirectMessage{ClientID: clientID, Message: message}
	}
	return nil
}

// Close closes the WebSocket handler and message broker
func (h *WSHandler) Close() error {
	return h.broker.Close()
}

// GetConfig returns the WebSocket configuration
func (h *WSHandler) GetConfig() WSConfig {
	return h.config
}

// GetBroker returns the message broker
func (h *WSHandler) GetBroker() MessageBroker {
	return h.broker
}
