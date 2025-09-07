package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gurkanindibay/udemy-rest-api/models"
)

// getEvents godoc
// @Summary Get all events
// @Description Retrieve a list of all events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Failure 500 {object} map[string]string
// @Router /events [get]
func getEvents(c *gin.Context) {
	events, err := eventService.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// getEventByID godoc
// @Summary Get event by ID
// @Description Retrieve a specific event by its ID
// @Tags events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} models.Event
// @Failure 500 {object} map[string]string
// @Router /events/{id} [get]
func getEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := eventService.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// createEvent godoc
// @Summary Create a new event
// @Description Create a new event (requires authentication)
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event body models.CreateEventRequest true "Event data"
// @Success 201 {object} models.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
// @Security BearerAuth
func createEvent(c *gin.Context) {

	userId := c.GetInt64("userId")

	var request models.CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEvent := models.Event{
		Name:        request.Name,
		Description: request.Description,
		Location:    request.Location,
		DateTime:    request.DateTime,
		UserId:      userId,
	}

	createdEvent, err := eventService.CreateEvent(newEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdEvent)
}

// updateEvent godoc
// @Summary Update an event
// @Description Update an existing event (requires authentication and ownership)
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Param event body models.CreateEventRequest true "Updated event data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id} [put]
// @Security BearerAuth
func updateEvent(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetInt64("userId")

	// Check if the event exists and belongs to the user
	event, err := eventService.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this event"})
		return
	}

	// Bind the updated event data
	var request models.CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert id from string to int64
	eventID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	updatedEvent := models.Event{
		ID:          eventID,
		Name:        request.Name,
		Description: request.Description,
		Location:    request.Location,
		DateTime:    request.DateTime,
		UserId:      userId,
	}
	if err := eventService.UpdateEvent(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// deleteEvent godoc
// @Summary Delete an event
// @Description Delete an existing event (requires authentication and ownership)
// @Tags events
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Event ID"
// @Success 204
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id} [delete]
// @Security BearerAuth
func deleteEvent(c *gin.Context) {
	id := c.Param("id")

	userId := c.GetInt64("userId")

	// Check if the event exists and belongs to the user
	event, err := eventService.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this event"})
		return
	}

	if err := eventService.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
