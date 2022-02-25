package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-house-plant-tracker-api/helpers"
)

type Plant struct {
	Name    string `bson:"name"`
	Type    string `bson:"type,omitempty"`
	Species string `bson:"species,omitempty"`
	Image   string `bson:"image,omitempty"`
	Active  bool   `bson:"active"`
}

func SavePlant(ginContext *gin.Context) {
	var input Plant

	client, context := helpers.Connect()
	collection := client.Database(helpers.GetSetting("DATABASE.NAME")).Collection("plants")

	if err := ginContext.ShouldBindJSON(&input); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newData := Plant{
		Name:    input.Name,
		Type:    input.Type,
		Species: input.Species,
		Image:   input.Image,
		Active:  true,
	}

	result, err := collection.InsertOne(context, newData)

	if err != nil {
		panic(err)
	}

	ginContext.JSON(http.StatusOK, result.InsertedID)

	client.Disconnect(context)
}
