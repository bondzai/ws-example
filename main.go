package main

import (
	"api-gateway/config"
	"api-gateway/internal/graph"
	"api-gateway/internal/handlers"
	"api-gateway/internal/infrastructures"
	"api-gateway/internal/repositories"
	"api-gateway/internal/servers"
	"api-gateway/internal/usecases"
)

func main() {
	conf := config.NewConfig()

	// Infrastructures
	db, dbConn := infrastructures.NewGorm(conf.DbDsn)
	defer dbConn.Close()

	redisClient := infrastructures.NewRedisFromUrl(conf.RedisUrl)
	defer redisClient.Conn().Close()

	// Repositories
	merchantRepository := repositories.NewMerchantRepository(db)
	merchantCategoryRepository := repositories.NewMerchantCategoryRepository(db)
	cryptoChainRepository := repositories.NewCryptoChainRepository(db)
	userRepository := repositories.NewUserRepository()
	merchantCacheRepository := repositories.NewMerchantCacheRepository(redisClient)

	// UseCases
	merchantUseCase := usecases.NewMerchantUseCase(merchantRepository, merchantCacheRepository)
	merchantCategoryUseCase := usecases.NewMerchantCategoryUseCase(merchantCategoryRepository)
	cryptoChainUseCase := usecases.NewCryptoChainUseCase(cryptoChainRepository)
	userUseCase := usecases.NewUserUseCase(userRepository)

	// Handlers
	graphqlHandler := handlers.NewGraphqlHandler(&graph.Resolver{
		MerchantUseCase: merchantUseCase,
	})
	userWsHandler := handlers.NewUserWebsocketHandler(userUseCase)
	restHandler := handlers.NewRestHandler(
		merchantUseCase,
		merchantCategoryUseCase,
		cryptoChainUseCase,
	)

	// Start HTTP Server
	servers.StartHttpServer(
		conf,
		graphqlHandler,
		restHandler,
		userWsHandler,
	)
}
