package ws

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

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
			for client := range handler.hub.Clients {
				if client.ID == msg.ClientID {
					client.SendMessage(msg.Data)
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

// Broadcast sends a message to all connected clients.
func (h *WSHandler) Broadcast(message []byte) {
	h.hub.Broadcast <- message
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
		// Fallback: locally broadcast to all clients with the same ID, excluding sender.
		h.hub.mu.Lock()
		defer h.hub.mu.Unlock()
		for client := range h.hub.Clients {
			if client.ID == sender.ID && client != sender {
				client.SendMessage(message)
			}
		}
	}
}
