package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dbURL := os.Getenv("DATABASE")
	fmt.Println(dbURL)
	opts, e := pg.ParseURL(dbURL)
	if e != nil {
		panic(e)
	}
	db := pg.Connect(opts)
	defer db.Close()
	if _, e := db.Exec("SELECT 1"); e != nil {
		panic(e)
	}
	fmt.Println("Database connection established")
	e = createSchema(db)
	if e != nil {
		panic(e)
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

		user := &User{
			Username: body.Username,
			Password: body.Password,
		}
		if e := db.Insert(user); e != nil {
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID       int64
	Username string `pg:",pk,notnull,unique"`
	Password string `pg:",notnull"`
}

func (u *User) String() string {
	return fmt.Sprintf("User<%d, %s>", u.ID, u.Username)
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
	}
	for _, model := range models {
		e := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if e != nil {
			return e
		}
	}
	return nil
}
