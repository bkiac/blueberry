package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bkiac/blueberry/server/ent"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE")
	db, e := ent.Open("postgres", dsn)
	if e != nil {
		log.Fatalf("Failed to open database: %v", e)
		panic(e)
	}
	defer db.Close()
	// Migration
	if e := db.Schema.Create(context.Background()); e != nil {
		log.Fatalf("Failed to create schema resources: %v", e)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	r.POST("/user", func(c *gin.Context) {
		var body Register
		if e := c.ShouldBindJSON(&body); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			return
		}

		user, e := CreateUser(c, db, body)
		if e != nil {
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Register is the validator for registration POST request JSON body.
type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser creates a user in the database.
func CreateUser(c context.Context, db *ent.Client, args Register) (*ent.User, error) {
	u, e := db.User.Create().SetUsername(args.Username).SetPassword(args.Password).Save(c)
	if e != nil {
		return nil, e
	}
	return u, nil
}
