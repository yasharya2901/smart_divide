package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(rg *gin.RouterGroup) {
	expenses := rg.Group("/expenses")

	expenses.GET("/", func(c *gin.Context) {
		// For now, mock a response
		c.JSON(http.StatusOK, gin.H{"message": "Fetching all expenses"})
	})

	expenses.POST("/", func(c *gin.Context) {
		// Mock response
		c.JSON(http.StatusCreated, gin.H{"message": "Expense added successfully"})
	})
}
