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
	"github.com/gorilla/websocket"
	"github.com/tomato3713/storyline/server/graph"
	"github.com/tomato3713/storyline/server/graph/resolver"
	"github.com/tomato3713/storyline/server/repository"
	"github.com/tomato3713/storyline/server/services"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()
	client, db, err := connectDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	repo := repository.NewRepository(db)

	services := services.New(repo)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &resolver.Resolver{
					Srv: services,
				},
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Host == "localhost:3000"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	corsHandler := CORSMiddleware(srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", corsHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectDatabase(ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	// MongoDBの接続URI
	uri := "mongodb://admin:password@localhost:27017"

	// MongoDBクライアントのオプションを設定
	clientOptions := options.Client().ApplyURI(uri)

	// クライアントを作成
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create MongoDB client: %w", err)
	}

	// 接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB!")

	database := client.Database("storyline")
	return client, database, nil
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
