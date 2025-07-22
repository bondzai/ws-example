package main

import (
	"context"
	"log"
	"path/filepath"

	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/infrastructures"
	"api-gateway/internal/repositories"
	"api-gateway/internal/usecases"
	"api-gateway/pkg/filestorage"
	"api-gateway/pkg/ws"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	conf := config.NewConfig()

	redisClient := infrastructures.NewRedis(
		conf.Redis.URI,
		conf.Redis.Password,
		conf.Redis.DB,
	)

	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(conf.Mongo.URI),
	)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	mongoDB := mongoClient.Database(conf.Mongo.Database)

	// --- Repositories ---
	userRepository := repositories.NewMockUserRepository()
	messageRepository := repositories.NewMongoMessageRepository(mongoDB)

	// --- File Storage ---
	uploadsPath, _ := filepath.Abs("./uploads")
	fileStorage, err := filestorage.NewLocalStorage(uploadsPath, conf.BaseUrl+"/files")
	if err != nil {
		log.Fatalf("Failed to create file storage: %v", err)
	}

	// --- WebSockets ---
	connManager := ws.NewConnectionManager(
		ws.WithRedis(redisClient),
		ws.WithAutoSync(true),
	)

	// --- Use Cases ---
	chatUseCase := usecases.NewChatUseCase(userRepository, messageRepository, connManager)
	fileUploadUseCase := usecases.NewFileUploadUseCase(fileStorage)

	// --- Handlers ---
	chatHandler := handlers.NewChatHandler(chatUseCase, connManager)
	fileUploadHandler := handlers.NewFileUploadHandler(fileUploadUseCase)

	app := infrastructures.NewFiber()

	// --- Static File Server ---
	app.Static("/files", uploadsPath)

	// --- API Routes ---
	v1 := app.Group("/api/v1")
	{
		wsGroup := v1.Group("/ws")
		wsGroup.Get("/chat", chatHandler.ServeWS)

		fileGroup := v1.Group("/files")
		fileGroup.Post("/upload", fileUploadHandler.UploadFile)
	}

	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Fatal(app.Listen(":" + conf.HttpPort))
}
