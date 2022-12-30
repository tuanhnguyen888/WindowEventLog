package main

import (
	"github.com/gin-gonic/gin"
)

const port = "81"

func main() {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	app.Use(CORSMiddleware())

	app.POST("/", Broker)
	app.POST("/handle", Handle)

	err := app.Run(":" + port)
	if err != nil {
		panic("can not run port ")
	}
}
