package handlers

import (
	"api-gateway/pkg/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// NewChatHandler creates a new WebSocket handler for the chat.
// It takes the core ws.WSHandler and returns a Fiber-compatible handler.
func NewChatHandler(wsHandler *ws.WSHandler) fiber.Handler {
	return websocket.New(wsHandler.Handler)
}
