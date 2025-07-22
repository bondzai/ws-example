package main

import (
	"context"
	"log"

	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/infrastructures"
	"api-gateway/internal/repositories"
	"api-gateway/internal/usecases"
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

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	mongoDB := mongoClient.Database(conf.Mongo.Database)

	userRepository := repositories.NewMockUserRepository()
	messageRepository := repositories.NewMongoMessageRepository(mongoDB)

	wsHandler := ws.NewWSHandler(
		ws.WithRedis(redisClient),
		ws.WithAutoSync(true),
	)

	_ = usecases.NewChatUseCase(userRepository, messageRepository, wsHandler)

	chatHandler := handlers.NewChatHandler(wsHandler)

	app := infrastructures.NewFiber()

	v1 := app.Group("/api/v1")
	{
		wsGroup := v1.Group("/ws")
		wsGroup.Get("/chat", chatHandler)
	}

	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Fatal(app.Listen(":" + conf.HttpPort))
}
