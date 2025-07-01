package infrastructures

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiber() *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
	}))

	return app
}
