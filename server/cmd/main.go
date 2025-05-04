package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tomato3713/storyline/server/graph"
	"github.com/tomato3713/storyline/server/repository"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()
	db, err := connectDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Repository: repo,
				},
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectDatabase(ctx context.Context) (*mongo.Database, error) {
	// MongoDBの接続URI
	uri := "mongodb://root:example@localhost:27017"

	// MongoDBクライアントのオプションを設定
	clientOptions := options.Client().ApplyURI(uri)

	// クライアントを作成
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to create MongoDB client: %w", err)
	}
	defer client.Disconnect(ctx)

	// 接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB!")

	database := client.Database("storyline")
	return database, nil
}
