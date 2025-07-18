package ws

import (
	"context"
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
						log.Printf("[DEBUG] Room %s deleted (no clients left).", client.RoomID)
					}
				}
			}
			h.mu.Unlock()

		case roomMsg := <-h.Broadcast:
			h.mu.Lock()
			room, ok := h.Rooms[roomMsg.RoomID]
			if ok {
				clientIDs := []string{}
				for client := range room.Clients {
					clientIDs = append(clientIDs, client.ID)
				}
				log.Printf("[DEBUG] Broadcasting to room %s. Clients: %v. Message: %s", roomMsg.RoomID, clientIDs, string(roomMsg.Message))
				for client := range room.Clients {
					select {
					case client.Send <- roomMsg.Message:
					default:
						close(client.Send)
						delete(room.Clients, client)
					}
				}
			} else {
				log.Printf("[DEBUG] Tried to broadcast to non-existent room %s.", roomMsg.RoomID)
			}
			h.mu.Unlock()
		}
	}
}

// WSHandler provides customizable callbacks and a hub for WebSocket events.
type WSHandler struct {
	hub          *Hub
	OnConnect    func(c *Client)
	OnMessage    func(c *Client, messageType int, data []byte)
	OnClose      func(c *Client)
	PingInterval time.Duration
	// Sync is the sync adapter for cross-node message propagation.
	// It is not coupled with any specific implementation.
	Sync SyncAdapter
	// AutoSync, when true, automatically syncs received messages to other clients with the same ID.
	AutoSync bool
}

// SyncMessage defines the structure for synchronization messages.
type SyncMessage struct {
	ClientID string
	Data     []byte
}

// SyncAdapter abstracts the pub/sub system used to sync messages across nodes.
type SyncAdapter interface {
	// Publish sends a sync message to the backend.
	Publish(ctx context.Context, msg SyncMessage) error
	// Subscribe starts listening for sync messages and calls handler on each message.
	Subscribe(ctx context.Context, handler func(msg SyncMessage))
}

// NewWSHandler initializes a new WSHandler with its hub.
// If a Sync adapter is provided, it subscribes to incoming sync messages.
func NewWSHandler() *WSHandler {
	handler := &WSHandler{
		hub:          NewHub(),
		PingInterval: 30 * time.Second,
		AutoSync:     false, // default disabled; enable it for auto sync behavior
	}
	go handler.hub.Run()

	// If a Sync adapter is set, subscribe to sync messages.
	if handler.Sync != nil {
		handler.Sync.Subscribe(context.Background(), func(msg SyncMessage) {
			handler.hub.mu.Lock()
			defer handler.hub.mu.Unlock()
			for _, room := range handler.hub.Rooms {
				for client := range room.Clients {
					if client.ID == msg.ClientID {
						client.SendMessage(msg.Data)
					}
				}
			}
		})
	}

	return handler
}

// Handler returns a Fiber-compatible WebSocket handler function.
func (h *WSHandler) Handler(c *websocket.Conn) {
	client := NewClient(c, h)
	h.hub.Register <- client
	if h.OnConnect != nil {
		h.OnConnect(client)
	}
	client.Listen() // Blocks until the connection closes.
}

// BroadcastToRoom sends a message to all clients in a specific room.
func (h *WSHandler) BroadcastToRoom(roomID string, message []byte) {
	h.hub.Broadcast <- RoomMessage{RoomID: roomID, Message: message}
}

// Broadcast sends a message to all connected clients.
func (h *WSHandler) Broadcast(message []byte) {
	h.hub.Broadcast <- RoomMessage{RoomID: "all", Message: message}
}

// SyncMessage publishes a synchronization message using the provided Sync adapter.
// If a Sync adapter is not configured, it falls back to locally forwarding the message
// to all other clients with the same ID.
func (h *WSHandler) SyncMessage(ctx context.Context, sender *Client, message []byte) {
	if h.Sync != nil {
		msg := SyncMessage{
			ClientID: sender.ID,
			Data:     message,
		}
		if err := h.Sync.Publish(ctx, msg); err != nil {
			log.Println("failed to publish sync message:", err)
		}
	} else {
		// Fallback: locally broadcast to all clients with the same ID in the same room, excluding sender.
		h.hub.mu.Lock()
		defer h.hub.mu.Unlock()
		room, ok := h.hub.Rooms[sender.RoomID]
		if ok {
			for client := range room.Clients {
				if client.ID == sender.ID && client != sender {
					client.SendMessage(message)
				}
			}
		}
	}
}
