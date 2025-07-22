package handlers

import (
	"api-gateway/internal/usecases"
	"api-gateway/pkg/ws"
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// ChatHandler handles the WebSocket connections for the chat.
type ChatHandler struct {
	useCase     usecases.ChatUseCase
	connManager *ws.ConnectionManager
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(useCase usecases.ChatUseCase, connManager *ws.ConnectionManager) *ChatHandler {
	return &ChatHandler{
		useCase:     useCase,
		connManager: connManager,
	}
}

// ServeWS is the entry point for WebSocket connections.
func (h *ChatHandler) ServeWS(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		// Create a new client from the WebSocket connection.
		client := ws.NewClient(conn, h.connManager)
		h.connManager.RegisterClient(client)
		defer h.connManager.UnregisterClient(client)

		// --- OnConnect ---
		// Notify the use case that a user has connected and get chat history.
		history, err := h.useCase.UserConnected(context.Background(), client.GetID(), client.GetRoomID())
		if err != nil {
			log.Printf("Error on user connected: %v", err)
			conn.Close() // Close connection if setup fails.
			return
		}

		// Send chat history to the newly connected client.
		for _, msg := range history {
			payload, _ := json.Marshal(msg)
			client.SendMessage(payload)
		}

		// Start a listening goroutine for the client.
		go client.WritePump()

		// --- OnMessage ---
		// Read messages from the client in a loop.
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// --- OnDisconnect ---
				log.Printf("Client %s disconnected: %v", client.GetID(), err)
				h.useCase.UserDisconnected(context.Background(), client.GetID(), client.GetRoomID())
				break
			}
			// Process the message using the use case.
			h.useCase.ProcessMessage(context.Background(), client.GetID(), client.GetRoomID(), string(msg))
		}
	})(c)
}
