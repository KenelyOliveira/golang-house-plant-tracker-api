package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlantType struct {
	Title  string `bson:"title,omitempty"`
	Active bool   `bson:"active,omitempty"`
}

func getEnvVar(name string) string {
	data := os.Getenv(name)

	if data == "" {
		log.Fatal(fmt.Sprintf("$%s must be set", name))
		return ""
	}

	return data
}

func Connect() (*mongo.Client, context.Context) {
	cnnStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		getEnvVar("DATABASE.USER"),
		getEnvVar("DATABASE.PASSWORD"),
		getEnvVar("DATABASE.URI"),
		getEnvVar("DATABASE.NAME"))

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

	collection := client.Database(getEnvVar("DATABASE.NAME")).Collection("types")
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

func main() {
	port := getEnvVar("PORT")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/plant-types", getPlantTypes)

	router.Run(":" + port)
}
