package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-house-plant-tracker-api/controllers"
	"github.com/golang-house-plant-tracker-api/helpers"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := helpers.GetSetting("PORT")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/plant-types", controllers.GetPlantTypes)
	router.POST("/plant", controllers.SavePlant)
	//router.POST("/plant", uploadPicture)

	router.Run(":" + port)
}
