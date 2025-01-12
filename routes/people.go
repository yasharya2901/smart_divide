package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/handlers"
	"gorm.io/gorm"
)

func PersonRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Group for people-related routes
	people := rg.Group("/people")
	peopleHandler := handlers.NewPeopleHandler(db)

	people.GET("/:id", peopleHandler.GetPerson())
	people.PUT("/:id", peopleHandler.UpdatePerson())

	people.POST("/contacts", peopleHandler.GetPeopleByContacts())
	people.POST("/emails", peopleHandler.GetPeopleByEmails())
}
