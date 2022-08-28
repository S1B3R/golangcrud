package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func getUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	var filter, option interface{}

	filter = bson.D{
		{Key: "_id", Value: id},
	}
	option = bson.D{}

	cursor, err := query(mongoClient, c, dbName, usrCol, filter, option)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	var results []User
	if err := cursor.All(c, &results); err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	if len(results) <= 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.IndentedJSON(http.StatusOK, results[len(results)-1])
}

func getUsers(c *gin.Context) {
	var filter, option interface{}

	filter = bson.D{}
	option = bson.D{}

	cursor, err := query(mongoClient, c, dbName, usrCol, filter, option)
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

	result, err := insertOne(mongoClient, c, dbName, usrCol, newUser)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	c.IndentedJSON(http.StatusOK, result)
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

	var results []User
	cursor, err := query(mongoClient, c, dbName, usrCol, filter, option)
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

		_, err := ReplaceOne(mongoClient, c, dbName, usrCol, filter, value)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}
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
	result, err = deleteOne(mongoClient, c, dbName, usrCol, filter)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}

	if result.DeletedCount <= 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

}
