package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bkiac/blueberry/server/database"
	"github.com/bkiac/blueberry/server/graphql"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	var err error

	database.SetupDatabase()
	graphql.SetupSchema()

	// Starts GraphQL server.
	http.Handle("/", &relay.Handler{Schema: Schema})
	log.Println("📡 GraphQL server is listening!")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}
}
