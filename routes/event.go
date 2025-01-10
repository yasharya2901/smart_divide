package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/controllers"
	"gorm.io/gorm"
)

func EventRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	event := rg.Group("/events")

	var eventController = controllers.NewEventController(db)

	event.GET("/", func(c *gin.Context) {

		events, err := eventController.GetEvents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"events": events})
	})

	event.GET("/:id", func(c *gin.Context) {
		idAsString := c.Param("id")
		// convert id to uint
		id, err := strconv.ParseUint(idAsString, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		preloadPerson := c.Query("preloadPerson") == "true"
		event, err := eventController.GetEventByID(uint(id), preloadPerson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"event": event})
	})

	event.POST("/", func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		event, err := eventController.CreateEvent(input.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"event": event})
	})
}
