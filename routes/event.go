package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/handlers"
	"gorm.io/gorm"
)

func EventRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	event := rg.Group("/events")
	var eventHandler = handlers.NewEventHandler(db)

	// Base event routes
	event.GET("/", eventHandler.GetEvents())
	event.POST("/", eventHandler.CreateEvent())

	// Single event routes
	event.GET("/:id", eventHandler.GetEvent())
	event.PUT("/:id", eventHandler.UpdateEvent())
	event.DELETE("/:id", eventHandler.DeleteEvent())

	// People management routes - use different base path
	people := event.Group("/:id/members")
	people.POST("/:personId", eventHandler.AddPersonToEvent())
	people.DELETE("/:personId", eventHandler.RemovePersonFromEvent())
}
