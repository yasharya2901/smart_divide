package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/services"
	"gorm.io/gorm"
)

type PeopleHandler struct {
	service *services.PeopleService
}

func NewPeopleHandler(db *gorm.DB) *PeopleHandler {
	return &PeopleHandler{service: services.NewPeopleService(db)}
}

type response struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
}

func (p *PeopleHandler) CreatePerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name    string `json:"name" binding:"required"`
			Contact string `json:"contact" binding:"required"`
			Email   string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		person, err := p.service.CreatePerson(req.Name, req.Contact, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"person": response{person.ID, person.Name, person.Contact, person.Email}})
	}
}

func (p *PeopleHandler) GetPerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		personID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		person, err := p.service.GetPersonByID(uint(personID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"people": person})
	}
}

func (p *PeopleHandler) UpdatePerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		peopleID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var req struct {
			Name    string `json:"name"`
			Contact string `json:"contact"`
			Email   string `json:"email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		person, err := p.service.UpdatePerson(uint(peopleID), req.Name, req.Contact, req.Email, "", time.Time{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"person": response{ID: person.ID, Name: person.Name, Contact: person.Contact, Email: person.Email}})

	}
}

func (p *PeopleHandler) GetPeopleByEmails() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Emails []string `json:"emails" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		people, err := p.service.GetPeopleByEmails(request.Emails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := make([]response, len(people))

		for i, person := range people {
			res[i] = response{
				ID:      person.ID,
				Name:    person.Name,
				Contact: person.Contact,
				Email:   person.Email,
			}
		}

		c.JSON(http.StatusOK, gin.H{"people": res})
	}
}

func (p *PeopleHandler) GetPeopleByContacts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Contacts []string `json:"contacts" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		people, err := p.service.GetPeopleByContacts(request.Contacts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := make([]response, len(people))

		for i, person := range people {
			res[i] = response{
				ID:      person.ID,
				Name:    person.Name,
				Contact: person.Contact,
				Email:   person.Email,
			}
		}

		c.JSON(http.StatusOK, gin.H{"people": res})
	}
}
