package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dbURL := os.Getenv("DATABASE")
	fmt.Println(dbURL)
	opts, err := pg.ParseURL(dbURL)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opts)
	defer db.Close()
	if _, err := db.Exec("SELECT 1"); err != nil {
		panic(err)
	}
	fmt.Println("Database connection established")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
