package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"saveEventLog/database"
)

var port = "82"

func main() {

	redisClient, err := database.NewInitRedis()
	if err != nil {
		log.Panic(err)
	}

	pgDb, err := database.NewInitPG()
	if err != nil {
		log.Panic(err)
	}

	saveLogHandler := NewSaveLogHandler(redisClient, pgDb)

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())
	app.Use(CORSMiddleware())

	app.GET(
		"/savePg",
		saveLogHandler.saveLogByPg,
	)

	app.GET(
		"/find",
		saveLogHandler.showLogFromPG,
	)

	err = app.Run(":" + port)
	if err != nil {
		panic("can not run port ")
	}

}
