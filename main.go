package main

import (
	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/infrastructures"
	"api-gateway/internal/repositories"
	"api-gateway/internal/usecases"
	"log"
)

func main() {
	// Load application configuration
	conf := config.NewConfig()

	// Initialize dependencies
	redisClient := infrastructures.NewRedis(conf.Redis.URL, conf.Redis.Password, conf.Redis.DB)
	userRepository := repositories.NewUserRepository()
	userUseCase := usecases.NewUserUseCase(userRepository)

	// Inject dependencies into the handler
	userWsHandler := handlers.NewUserWebsocketHandler(userUseCase, redisClient)

	// Set up the Fiber application
	app := infrastructures.NewFiber()

	v1 := app.Group("/api/v1")
	{
		ws := v1.Group("/ws")
		ws.Get("/chat", userWsHandler)
	}

	// Start the server
	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Fatal(app.Listen(":" + conf.HttpPort))
}
