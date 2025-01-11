package handlers

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/services"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(db *gorm.DB) *EventHandler {
	return &EventHandler{service: services.NewEventService(db)}
}

func (h *EventHandler) AddPersonToEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}

		personID, err := strconv.ParseUint(c.Param("personId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
			return
		}

		if err := h.service.AddPersonToEvent(uint(eventID), uint(personID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (h *EventHandler) GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := h.service.GetEvents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"events": events})
	}
}

func (h *EventHandler) GetEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		idAsString := c.Param("id")
		id, err := strconv.ParseUint(idAsString, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		preloadPerson := c.Query("preloadPerson") == "true"
		event, err := h.service.GetEventByID(uint(id), preloadPerson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"event": event})
	}
}

func (h *EventHandler) CreateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		event, err := h.service.CreateEvent(input.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"event": event})
	}
}

func (h *EventHandler) UpdateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		idAsString := c.Param("id")
		id, err := strconv.ParseUint(idAsString, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name == "" {
			c.JSON(http.StatusNoContent, gin.H{"error": "Name is required"})
			return
		}

		event, err := h.service.UpdateEvent(uint(id), input.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"event": event})
	}
}

func (h *EventHandler) DeleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		idAsString := c.Param("id")
		id, err := strconv.ParseUint(idAsString, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.DeleteEvent(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *EventHandler) RemovePersonFromEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}

		personID, err := strconv.ParseUint(c.Param("personId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
			return
		}

		if err := h.service.RemovePersonFromEvent(uint(eventID), uint(personID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
