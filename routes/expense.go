package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/handlers"
	"gorm.io/gorm"
)

func ExpenseRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	expenses := rg.Group("/expenses")

	expenseHandler := handlers.NewExpenseHandler(db)

	// Base expense routes
	expenses.GET("/", expenseHandler.GetExpenses())
	expenses.POST("/", expenseHandler.CreateExpense())

	// Single expense routes
	expenses.GET("/:id", expenseHandler.GetExpense())
	expenses.PUT("/:id", expenseHandler.UpdateExpense())
	expenses.DELETE("/:id", expenseHandler.DeleteExpense())

	// Expense participants management routes
	participants := expenses.Group("/:id/participants")
	participants.GET("/", expenseHandler.GetParticipants())
	participants.POST("/", expenseHandler.AddParticipant())
	participants.PUT("/:personId", expenseHandler.UpdateParticipant())
	participants.DELETE("/:personId", expenseHandler.RemoveParticipant())

	// Check for payment consistency
	expenses.GET("/:id/check", expenseHandler.CheckExpenseConsistency())

}
