package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/models"
	"github.com/gurkanindibay/udemy-rest-api/db"
)

func main() {
	db.InitDB("events.db")
	server := gin.Default()
	server.GET("/events", getEvents)
	server.POST("/events", createEvent)
	server.Run(":8080")
}

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func createEvent(c *gin.Context) {
	var newEvent models.Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newEvent.ID = 1
	newEvent.UserId = 1 // later: get from auth context
	if err := newEvent.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newEvent)
}
