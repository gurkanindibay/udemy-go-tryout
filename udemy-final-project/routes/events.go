package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/models"
	"github.com/gurkanindibay/udemy-rest-api/utils"
)

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

func createEvent(c *gin.Context) {

	// validate JWT token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

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

func updateEvent(c *gin.Context) {
	id := c.Param("id")
	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert id from string to int64
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	updatedEvent.ID = eventID
	if err := updatedEvent.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

func deleteEvent(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
