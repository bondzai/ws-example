package handlers

import (
	"api-gateway/internal/graph"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/graphql-go/handler"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
)

type GraphqlHandler struct {
	handler *handler.Handler
}

func NewGraphqlHandler(resolver *graph.Resolver) *graphqlHandler.Server {
	gh := graphqlHandler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolver,
			},
		),
	)

	gh.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	gh.AddTransport(transport.Options{})
	gh.AddTransport(transport.GET{})
	gh.AddTransport(transport.POST{})
	gh.AddTransport(transport.MultipartForm{})
	gh.AddTransport(transport.UrlEncodedForm{})
	gh.AddTransport(transport.GRAPHQL{})

	gh.Use(extension.Introspection{})

	return gh
}

func (h *GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
