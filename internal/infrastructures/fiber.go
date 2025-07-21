package infrastructures

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  2 * time.Minute, // Increased for WebSocket connections
		WriteTimeout: 2 * time.Minute, // Increased for WebSocket connections
		IdleTimeout:  5 * time.Minute, // How long to keep idle connections open
		BodyLimit:    4 * 1024 * 1024, // 4MB body limit
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
	}))

	return app
}
