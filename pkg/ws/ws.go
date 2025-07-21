package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

// Room holds a set of clients in a chat room.
type Room struct {
	Clients map[*Client]bool
}

// Hub maintains the set of active rooms and broadcasts messages to rooms.
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan RoomMessage
	mu         sync.Mutex
}

// RoomMessage is a message to be broadcast to a specific room.
type RoomMessage struct {
	RoomID  string
	Message []byte
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan RoomMessage),
	}
}

// Run starts the hub loop to process client registration, unregistration, and broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			room, ok := h.Rooms[client.RoomID]
			if !ok {
				room = &Room{Clients: make(map[*Client]bool)}
				h.Rooms[client.RoomID] = room
			}
			room.Clients[client] = true
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
			h.mu.Unlock()

		case roomMsg := <-h.Broadcast:
			h.mu.Lock()
			room, ok := h.Rooms[roomMsg.RoomID]
			if ok {
				for client := range room.Clients {
					select {
					case client.Send <- roomMsg.Message:
					default:
						close(client.Send)
						delete(room.Clients, client)
					}
				}
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
	MessageBroker  MessageBrokerConfig
}

// WSHandler provides customizable callbacks and a hub for WebSocket events.
type WSHandler struct {
	hub       *Hub
	config    WSConfig
	broker    MessageBroker
	OnConnect func(c *Client)
	OnMessage func(c *Client, messageType int, data []byte)
	OnClose   func(c *Client)
}

// SyncMessage defines the structure for synchronization messages.
type SyncMessage struct {
	ClientID string `json:"client_id"`
	RoomID   string `json:"room_id"`
	Data     []byte `json:"data"`
	Type     string `json:"type"`
}

// NewWSHandler initializes a new WSHandler with its hub and message broker.
func NewWSHandler(opts ...Option) (*WSHandler, error) {
	// Start with the default configuration
	config := WSConfig{
		PingInterval:   30 * time.Second,
		PongWait:       60 * time.Second,
		WriteWait:      10 * time.Second,
		MaxMessageSize: 512,
		BufferSize:     256,
		EnableAutoSync: false,
		SyncChannel:    "websocket_sync",
		MessageBroker:  MessageBrokerConfig{Type: "noop"},
	}

	// Apply all the functional options to customize the configuration
	for _, opt := range opts {
		opt(&config)
	}

	// Create message broker
	brokerManager, err := NewMessageBrokerManager(config.MessageBroker)
	if err != nil {
		return nil, err
	}

	handler := &WSHandler{
		hub:    NewHub(),
		config: config,
		broker: brokerManager.GetBroker(),
	}

	go handler.hub.Run()

	// Subscribe to sync messages if auto-sync is enabled
	if config.EnableAutoSync {
		err = handler.broker.Subscribe(context.Background(), config.SyncChannel, handler.handleSyncMessage)
		if err != nil {
			log.Printf("Failed to subscribe to sync channel: %v", err)
		}
	}

	return handler, nil
}

// handleSyncMessage processes incoming sync messages from other nodes
func (h *WSHandler) handleSyncMessage(data []byte) {
	var syncMsg SyncMessage
	if err := json.Unmarshal(data, &syncMsg); err != nil {
		log.Printf("Failed to unmarshal sync message: %v", err)
		return
	}

	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()

	// Forward message to clients in the same room (excluding the sender)
	if room, ok := h.hub.Rooms[syncMsg.RoomID]; ok {
		for client := range room.Clients {
			if client.ID != syncMsg.ClientID {
				client.SendMessage(syncMsg.Data)
			}
		}
	}
}

// Handler returns a Fiber-compatible WebSocket handler function.
func (h *WSHandler) Handler(c *websocket.Conn) {
	client := NewClient(c, h)
	h.hub.Register <- client
	if h.OnConnect != nil {
		h.OnConnect(client)
	}

	client.Listen()
}

// BroadcastToRoom sends a message to all clients in a specific room.
func (h *WSHandler) BroadcastToRoom(roomID string, message []byte) {
	h.hub.Broadcast <- RoomMessage{RoomID: roomID, Message: message}
}

// Broadcast sends a message to all connected clients.
func (h *WSHandler) Broadcast(message []byte) {
	h.hub.Broadcast <- RoomMessage{RoomID: "all", Message: message}
}

// SyncMessage publishes a synchronization message using the configured message broker.
func (h *WSHandler) SyncMessage(ctx context.Context, sender *Client, message []byte) {
	if !h.config.EnableAutoSync {
		return
	}

	syncMsg := SyncMessage{
		ClientID: sender.ID,
		RoomID:   sender.RoomID,
		Data:     message,
		Type:     "message",
	}

	if err := h.broker.Publish(ctx, h.config.SyncChannel, syncMsg); err != nil {
		log.Printf("Failed to publish sync message: %v", err)
	}
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
