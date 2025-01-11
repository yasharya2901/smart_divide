package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/services"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{service: services.NewExpenseService(db)}
}

func (h *ExpenseHandler) GetExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {
		expenses, err := h.service.GetExpenses()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"expenses": expenses})
	}
}

func (h *ExpenseHandler) CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string  `json:"name" binding:"required"`
			TotalAmount float64 `json:"total_amount" binding:"required"`
			EventID     uint    `json:"event_id" binding:"required"`
			PaidByID    uint    `json:"paid_by_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense, err := h.service.CreateExpense(req.Name, req.TotalAmount, req.EventID, req.PaidByID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"expense": expense})
	}
}

func (h *ExpenseHandler) GetExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		uID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		expense, err := h.service.GetExpenseByID(uint(uID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"expense": expense})
	}
}

func (h *ExpenseHandler) UpdateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var req struct {
			Name        string  `json:"name"`
			TotalAmount float64 `json:"total_amount"`
			PaidByID    uint    `json:"paid_by_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense, err := h.service.UpdateExpense(uint(expenseID), req.Name, req.TotalAmount, req.PaidByID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"expense": expense})
	}
}

func (h *ExpenseHandler) DeleteExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.DeleteExpense(uint(expenseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *ExpenseHandler) GetParticipants() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		participants, err := h.service.GetExpensePeople(uint(expenseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"participants": participants})
	}
}

func (h *ExpenseHandler) AddParticipant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Expense ID"})
			return
		}

		var req struct {
			PersonID uint `json:"person_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		exp, err := h.service.AddExpensePerson(uint(expenseID), uint(req.PersonID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"participant": exp})
	}
}

func (h *ExpenseHandler) UpdateParticipant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Expense ID"})
			return
		}

		var req struct {
			PaidAmount float64 `json:"paid_amount"`
			OwedAmount float64 `json:"owed_amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		participantID := c.Param("participant_id")

		pID, err := strconv.ParseUint(participantID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Participant ID"})
			return
		}

		exp, err := h.service.UpdateExpensePerson(uint(expenseID), uint(pID), req.PaidAmount, req.OwedAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"participant": exp})
	}
}

func (h *ExpenseHandler) RemoveParticipant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Expense ID"})
			return
		}

		participantID := c.Param("participant_id")

		pID, err := strconv.ParseUint(participantID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Participant ID"})
			return
		}

		err = h.service.DeleteExpensePerson(uint(expenseID), uint(pID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *ExpenseHandler) CheckExpenseConsistency() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Expense ID"})
			return
		}

		err = h.service.CheckExpenseConsistency(uint(expenseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "true", "message": "Expense is consistent"})
	}
}
