package main

import (

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/db"
	"github.com/gurkanindibay/udemy-rest-api/routes"
)

func main() {
	db.InitDB("events.db")
	server := gin.Default()
	routes.SetupRoutes(server)
	server.Run(":8080")
}

