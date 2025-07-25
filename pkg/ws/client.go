package ws

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	// ID is a unique identifier for the client (e.g., extracted from a query parameter).
	ID string
	// Conn is the underlying WebSocket connection.
	Conn *websocket.Conn
	// Send is a buffered channel of outbound messages.
	Send chan []byte
	// handler holds the parent ConnectionManager.
	handler *ConnectionManager
	RoomID  string // Add RoomID to track which room the client is in
}

// NewClient creates a new Client instance.
// It extracts the client ID from the query parameter "id".
func NewClient(conn *websocket.Conn, handler *ConnectionManager) *Client {
	clientID := conn.Query("userId")
	roomID := conn.Query("roomId")

	return &Client{
		ID:      clientID,
		Conn:    conn,
		Send:    make(chan []byte, handler.config.BufferSize),
		handler: handler,
		RoomID:  roomID, // Set RoomID
	}
}

// Listen starts the read and write pumps.
func (c *Client) Listen() {
	go c.WritePump()
}

// readPump continuously reads messages from the WebSocket connection.
func (c *Client) readPump() {
	defer func() {
		c.handler.UnregisterClient(c)
		c.Conn.Close()
	}()

	// Set read deadline
	c.Conn.SetReadLimit(c.handler.config.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.handler.config.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(c.handler.config.PongWait))
		return nil
	})

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
		// The message processing is now handled in the ChatHandler,
		// so we remove the call to OnMessage from here.
		// Instead, the handler's ServeWS method will read messages.
		log.Printf("Received message type %d from %s", messageType, c.ID)
		_ = message
	}
}

// WritePump sends messages from the send channel and periodically pings the client.
func (c *Client) WritePump() {
	pingInterval := c.handler.config.PingInterval
	if pingInterval == 0 {
		pingInterval = 30 * time.Second
	}
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.handler.config.WriteWait))
			if !ok {
				// The channel was closed.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("write error:", err)
				return
			}
		case <-ticker.C:
			// Send a ping message.
			c.Conn.SetWriteDeadline(time.Now().Add(c.handler.config.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				log.Println("ping error:", err)
				return
			}
		}
	}
}

// SendMessage writes a message directly to the client's send channel.
func (c *Client) SendMessage(message []byte) {
	select {
	case c.Send <- message:
	default:
		log.Println("send channel full, dropping message")
	}
}

// GetID returns the client ID
func (c *Client) GetID() string {
	return c.ID
}

// GetRoomID returns the client's room ID
func (c *Client) GetRoomID() string {
	return c.RoomID
}

// GetHandler returns the WebSocket handler
func (c *Client) GetHandler() *ConnectionManager {
	return c.handler
}
