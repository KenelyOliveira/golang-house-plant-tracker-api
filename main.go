package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlantType struct {
	Title  string `bson:"title,omitempty"`
	Active bool   `bson:"active,omitempty"`
}

func Connect() (*mongo.Client, context.Context) {
	cnnStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.url"),
		viper.GetString("database.name"))

	client, err := mongo.NewClient(options.Client().ApplyURI(cnnStr))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

func getPlantTypes(ginContext *gin.Context) {
	client, context := Connect()

	collection := client.Database(viper.GetString("database.name")).Collection("types")
	cur, currErr := collection.Find(context, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(context)

	var types []PlantType
	if err := cur.All(context, &types); err != nil {
		panic(err)
	}

	ginContext.IndentedJSON(http.StatusOK, types)

	client.Disconnect(context)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/plant-types", getPlantTypes)

	router.Run(":" + port)
}
