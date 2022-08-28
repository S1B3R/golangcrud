package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var users = []User{}

var dbName string
var usrCol = "user"

var mongoClient *mongo.Client

func main() {
	dbName = os.Getenv("DB_NAME")
	dbString := os.Getenv("MONGO_STR")
	client, ctx, cancel, err := connect(dbString)
	if err != nil {
		panic(err)
	}
	mongoClient = client
	defer close(client, ctx, cancel)

	router := gin.Default()
	router.Use(validateAPIKey())

	router.GET("/user", getUsers)
	router.GET("/user/:id", getUser)
	router.POST("/user", addUser)
	router.PUT("/user", updateUser)
	router.DELETE("/user/:id", deleteUser)

	router.Run("localhost:8080")
}

func validateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		PRODkey := os.Getenv("X-API-Key")
		APIKey := c.Request.Header.Get("X-API-Key")
		if APIKey != PRODkey {
			c.AbortWithStatus(401)
			return
		}
		return
	}
}
