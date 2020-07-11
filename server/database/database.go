package database

import (
	"context"
	"log"
	"os"

	"github.com/bkiac/blueberry/server/database/ent"
)

// Client is the opened database client.
var Client *ent.Client

// SetupDatabase initializes Client.
func SetupDatabase(c context.Context) {
	var err error

	// Sets up database connection.
	dsn := os.Getenv("DATABASE")
	Client, err = ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Println("💽 Database opened successfully!")
	defer Client.Close()
	if err := Client.Schema.Create(c); err != nil {
		log.Fatalf("Failed to create schema resources: %v", err)
	}
}
