package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/handlers"
	"gorm.io/gorm"
)

func AuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	auth := rg.Group("/")
	var authHandler = handlers.NewAuthHandler(db)

	// Login and Register
	auth.POST("/login", authHandler.Login())
	auth.POST("/register", authHandler.Register())
}
