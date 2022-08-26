package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

var users = []user{
	{ID: 1, Name: "Gabriel de Paula", Age: 29, Email: "gabriel@gmail.com", Password: "000000", Address: "teste street"},
}

func main() {
	router := gin.Default()
	router.Use(validateAPIKey())
	router.GET("/users", getUsers)
	router.Run("localhost:8080")
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
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
