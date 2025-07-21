package main

import (
	"api-gateway/config"
	"api-gateway/pkg/ws"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	conf := config.NewConfig()

	// Example 1: Basic WebSocket handler without message broker
	basicHandler, err := ws.NewWSHandler(
		ws.WithAutoSync(false), // Explicitly disable auto-sync
	)
	if err != nil {
		log.Fatal("Failed to create basic handler:", err)
	}

	// Example 2: WebSocket handler with Redis message broker
	redisConfig := ws.RedisConfig{
		URL:      conf.MessageBroker.Redis.URL,
		Password: conf.MessageBroker.Redis.Password,
		DB:       conf.MessageBroker.Redis.DB,
		PoolSize: conf.MessageBroker.Redis.PoolSize,
	}

	redisHandler, err := ws.NewWSHandler(
		ws.WithRedis(redisConfig),
		ws.WithAutoSync(true),
		ws.WithSyncChannel("chat_sync"),
		ws.WithPingInterval(30*time.Second),
		ws.WithPongWait(60*time.Second),
		ws.WithWriteWait(10*time.Second),
		ws.WithMaxMessageSize(1024),
		ws.WithBufferSize(512),
	)
	if err != nil {
		log.Fatal("Failed to create Redis handler:", err)
	}

	// Example 3: Custom configuration
	customHandler, err := ws.NewWSHandler(
		ws.WithPingInterval(20*time.Second),
		ws.WithPongWait(40*time.Second),
		ws.WithWriteWait(5*time.Second),
		ws.WithMaxMessageSize(2048),
		ws.WithBufferSize(1024),
		ws.WithAutoSync(true),
		ws.WithSyncChannel("custom_sync"),
		ws.WithRedis(redisConfig),
	)
	if err != nil {
		log.Fatal("Failed to create custom handler:", err)
	}

	// Set up event handlers
	setupEventHandlers(basicHandler)
	setupEventHandlers(redisHandler)
	setupEventHandlers(customHandler)

	// Create Fiber app
	app := fiber.New()

	// WebSocket routes
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Basic WebSocket endpoint
	app.Get("/ws/basic", websocket.New(basicHandler.Handler))

	// Redis WebSocket endpoint
	app.Get("/ws/redis", websocket.New(redisHandler.Handler))

	// Custom WebSocket endpoint
	app.Get("/ws/custom", websocket.New(customHandler.Handler))

	// Start server
	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Fatal(app.Listen(":" + conf.HttpPort))
}

func setupEventHandlers(handler *ws.WSHandler) {
	// OnConnect event
	handler.OnConnect = func(c *ws.Client) {
		log.Printf("Client %s connected to room %s", c.GetID(), c.GetRoomID())

		// Send welcome message
		welcomeMsg := []byte(`{"type": "welcome", "message": "Welcome to the chat!"}`)
		c.SendMessage(welcomeMsg)
	}

	// OnMessage event
	handler.OnMessage = func(c *ws.Client, messageType int, data []byte) {
		log.Printf("Received message from client %s: %s", c.GetID(), string(data))

		// Echo the message back to the client
		response := []byte(`{"type": "echo", "message": "` + string(data) + `"}`)
		c.SendMessage(response)

		// Broadcast to room
		handler.BroadcastToRoom(c.GetRoomID(), data)
	}

	// OnClose event
	handler.OnClose = func(c *ws.Client) {
		log.Printf("Client %s disconnected from room %s", c.GetID(), c.GetRoomID())
	}
}

// Example of how to add a new message broker in the future
func exampleFutureKafkaBroker() {
	// This is an example of how you would add Kafka support in the future

	/*
		// 1. Add Kafka configuration to types.go
		type KafkaConfig struct {
			Brokers  []string
			Topic    string
			GroupID  string
			Username string
			Password string
		}

		// 2. Implement KafkaMessageBroker in message_broker.go
		type KafkaMessageBroker struct {
			*BaseMessageBroker
			writer *kafka.Writer
			reader *kafka.Reader
		}

		// 3. Add Kafka support to the factory
		func (f *MessageBrokerFactory) CreateMessageBroker(config MessageBrokerConfig) (MessageBroker, error) {
			switch config.Type {
			case "redis":
				return NewRedisMessageBroker(config.Redis)
			case "kafka":
				return NewKafkaMessageBroker(config.Kafka)
			case "noop", "":
				return NewNoOpMessageBroker(), nil
			default:
				return nil, fmt.Errorf("unsupported message broker type: %s", config.Type)
			}
		}

		// 4. Add helper function
		func WithKafka(kafkaConfig KafkaConfig) func(*WSConfig) {
			return func(config *WSConfig) {
				config.EnableAutoSync = true
				config.MessageBroker = MessageBrokerConfig{
					Type:   "kafka",
					Kafka:  kafkaConfig,
				}
			}
		}

		// 5. Use it in your application
		kafkaConfig := ws.KafkaConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "websocket_sync",
			GroupID: "websocket_group",
		}

		handler, err := ws.CreateHandler(
			ws.WithKafka(kafkaConfig),
			ws.WithAutoSync(true),
		)
	*/
}
