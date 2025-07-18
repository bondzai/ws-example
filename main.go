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
	conf := config.NewConfig()

	// Repositories
	userRepository := repositories.NewUserRepository()

	// UseCases
	userUseCase := usecases.NewUserUseCase(userRepository)

	// Handlers
	userWsHandler := handlers.NewUserWebsocketHandler(userUseCase)

	app := infrastructures.NewFiber()

	v1 := app.Group("/api/v1")

	{
		ws := v1.Group("/ws")
		ws.Get("/chat", userWsHandler)
	}

	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Fatal(app.Listen(":" + conf.HttpPort))
}
