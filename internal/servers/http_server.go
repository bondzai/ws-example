package servers

import (
	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/infrastructures"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func StartHttpServer(
	conf *config.Config,
	graphqlHandler *handler.Server,
	restHandler handlers.RestHandler,
	userWsHandler fiber.Handler,
) {
	app := infrastructures.NewFiber()

	app.Get("/", restHandler.HealthCheck)

	v1 := app.Group("/api/v1")

	{
		graph := v1.Group("/graphql")
		graph.All("", adaptor.HTTPHandlerFunc(graphqlHandler.ServeHTTP))
		graph.Get("/playground", func(c *fiber.Ctx) error {
			c.Set("Content-Type", "text/html")
			return c.Send(getGraphqlSandboxHtml(conf.BaseUrl))
		})
	}

	{
		merchants := v1.Group("/merchants")
		merchants.Post("", restHandler.CreateMerchant)
		merchants.Get("", restHandler.GetMerchants)
		merchants.Get("/:id", restHandler.GetMerchantById)
		merchants.Patch("/:id", restHandler.UpdateMerchant)
		merchants.Delete("", restHandler.DeleteMerchant)

		merchantCategories := v1.Group("/merchants-categories")
		merchantCategories.Post("", restHandler.CreateMerchantCategory)
		merchantCategories.Get("", restHandler.GetMerchantCategories)
		merchantCategories.Get("/:id", restHandler.GetMerchantCategoryById)
		merchantCategories.Patch("/:id", restHandler.UpdateMerchantCategory)
		merchantCategories.Delete("/:id", restHandler.DeleteMerchantCategory)
	}

	{
		cryptoChains := v1.Group("/crypto-chains")
		cryptoChains.Post("", restHandler.CreateCryptoChain)
		cryptoChains.Get("", restHandler.GetCryptoChains)
		cryptoChains.Get("/:id", restHandler.GetCryptoChainById)
		cryptoChains.Patch("/:id", restHandler.UpdateCryptoChain)
	}

	{
		ws := v1.Group("/ws")
		ws.Get("/users", userWsHandler)
	}

	log.Printf("Server is running on port: %s", conf.HttpPort)
	log.Printf("Access Apollo Sandbox at %s/graphql/playground", getMainURL(conf.BaseUrl))
	log.Fatal(app.Listen(":" + conf.HttpPort))
}
