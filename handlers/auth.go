package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/services"
	"github.com/yasharya2901/smart_divide/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{service: services.NewAuthService(db)}
}

func (a *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := a.service.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

func (a *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Contact  string `json:"contact" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !utils.ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}

		if !utils.ValidatePhoneNumber(req.Contact) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
			return
		}

		accessToken, refreshToken, err := a.service.Register(req.Name, req.Contact, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}
