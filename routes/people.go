package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PersonRoutes(rg *gin.RouterGroup) {
	// Group for people-related routes
	people := rg.Group("/people")

	// Define a simple route to get all people
	people.GET("/", func(c *gin.Context) {
		// Mock response for now
		c.JSON(http.StatusOK, gin.H{"message": "Fetching all people"})
	})

	// Example: Add a new person
	people.POST("/", func(c *gin.Context) {
		// Mock response
		c.JSON(http.StatusCreated, gin.H{"message": "Person added successfully"})
	})
}
