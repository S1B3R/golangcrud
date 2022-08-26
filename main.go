package main

import (
	"net/http"

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
	router.GET("/users", getUsers)
	router.Run("localhost:8080")
}

func getUsers(c *gin.Context) {
	//c.IndentedJSON(http.StatusOK, users)
	c.JSON(http.StatusOK, users)
}
