package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var users = []User{}

var dbName string

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
	router.GET("/users", getUsers)
	router.GET("/user/:id", getUser)
	router.POST("/user", addUser)
	router.PUT("/user/:id", updateUser)
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

func getUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	var filter, option interface{}

	filter = bson.D{
		{Key: "_id", Value: id},
	}
	option = bson.D{}

	cursor, err := query(mongoClient, c, dbName, "user", filter, option)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	var results []User
	if err := cursor.All(c, &results); err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getUsers(c *gin.Context) {
	var filter, option interface{}

	filter = bson.D{}
	option = bson.D{}

	cursor, err := query(mongoClient, c, dbName, "user", filter, option)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	var results []User
	if err := cursor.All(c, &results); err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	c.IndentedJSON(http.StatusOK, results)
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashedPass)
	newUser.Id = primitive.NewObjectID()

	_, err := insertOne(mongoClient, c, dbName, "user", newUser)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}
}

func updateUser(c *gin.Context) {
	var user User
	var err error
	if err = c.BindJSON(&user); err != nil {
		return
	}

	var filter, option interface{}
	filter = bson.D{
		{Key: "_id", Value: user.Id},
	}
	option = bson.D{}

	var cursor *mongo.Cursor
	cursor, err = query(mongoClient, c, dbName, "user", filter, option)
	if err := c.BindJSON(&user); err != nil {
		return
	}

	var results []User
	if err := cursor.All(c, &results); err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	for _, value := range results {
		value.Name = user.Name
		value.Age = user.Age
		value.Email = user.Email
		value.Password = string(hashedPass)
		value.Address = user.Address
	}

}

func deleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	filter := bson.D{
		{Key: "_id", Value: id},
	}

	var result *mongo.DeleteResult
	result, err = deleteOne(mongoClient, c, dbName, "user", filter)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	if result.DeletedCount <= 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

}
