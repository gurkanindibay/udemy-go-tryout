package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// registerForEvent godoc
// @Summary Register for an event
// @Description Register the authenticated user for a specific event
// @Tags registrations
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Success 201 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id}/register [post]
// @Security BearerAuth
func registerForEvent(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetInt64("userId")

	event, err := eventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := eventService.RegisterForEvent(userID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully registered for the event"})
}

// getUserRegistrations godoc
// @Summary Get user registrations
// @Description Get all events that a user has registered for
// @Tags registrations
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {array} models.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/registrations [get]
// @Security BearerAuth
func getUserRegistrations(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	registrations, err := eventService.GetUserRegistrations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, registrations)
}

// cancelRegistration godoc
// @Summary Cancel event registration
// @Description Cancel the authenticated user's registration for a specific event
// @Tags registrations
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id}/register [delete]
// @Security BearerAuth
func cancelRegistration(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetInt64("userId")

	event, err := eventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := eventService.CancelRegistration(userID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully canceled registration for the event"})
}
