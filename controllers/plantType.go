package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-house-plant-tracker-api/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type PlantType struct {
	Title  string `bson:"title,omitempty"`
	Active bool   `bson:"active,omitempty"`
}

func GetPlantTypes(ginContext *gin.Context) {
	client, context := helpers.Connect()

	collection := client.Database(helpers.GetSetting("DATABASE.NAME")).Collection("types")
	cur, currErr := collection.Find(context, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(context)

	var types []PlantType
	if err := cur.All(context, &types); err != nil {
		panic(err)
	}

	ginContext.JSON(http.StatusOK, types)

	client.Disconnect(context)
}
