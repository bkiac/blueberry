package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	Username string `gorm:"not_null;unique;index:username"`
	Password string `gorm:"nut_null"`
}

func main() {
	dsn := os.Getenv("DATABASE")
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if e != nil {
		panic(e)
	}
	fmt.Println("Database connection established")

	db.AutoMigrate(&User{})

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

		user := &User{
			Username: body.Username,
			Password: body.Password,
		}
		db.Create(user)
		// if e := db.Create(user); e != nil {
		// 	c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		// 	return
		// }

		c.JSON(http.StatusOK, user)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
