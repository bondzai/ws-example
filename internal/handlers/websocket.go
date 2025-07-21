package handlers

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/usecases"
	"api-gateway/pkg/ws"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type userWebsocketHandler struct {
	handler *ws.WSHandler
	usecase usecases.UserUseCase
}

func NewUserWebsocketHandler(usecase usecases.UserUseCase) fiber.Handler {
	userHandler, err := ws.NewWSHandler(
		ws.WithPingInterval(30*time.Second),
		ws.WithPongWait(60*time.Second),
		ws.WithWriteWait(10*time.Second),
		ws.WithMaxMessageSize(512),
		ws.WithBufferSize(256),
		ws.WithAutoSync(false),
	)
	if err != nil {
		log.Fatalf("Failed to create websocket handler: %v", err)
	}

	h := &userWebsocketHandler{
		handler: userHandler,
		usecase: usecase,
	}

	userHandler.OnConnect = h.handleConnection
	userHandler.OnMessage = h.handleMessage
	userHandler.OnClose = h.handleDisconnection

	return websocket.New(h.handler.Handler)
}

func (h *userWebsocketHandler) handleConnection(c *ws.Client) {
	userCount := h.usecase.IncreaseRealtimeUser()
	response, _ := json.Marshal(entities.UserCountResponse{
		ActiveUsers: userCount,
	})

	h.handler.BroadcastToRoom(c.RoomID, response)
}

func (h *userWebsocketHandler) handleMessage(c *ws.Client, messageType int, data []byte) {
	log.Printf("user %s sent: %s", c.ID, data)
	h.handler.BroadcastToRoom(c.RoomID, data)
}

func (h *userWebsocketHandler) handleDisconnection(c *ws.Client) {
	userCount := h.usecase.DecreaseRealtimeUser()
	response, _ := json.Marshal(entities.UserCountResponse{
		ActiveUsers: userCount,
	})

	h.handler.BroadcastToRoom(c.RoomID, response)
}
