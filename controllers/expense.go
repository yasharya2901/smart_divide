package controllers

import (
	"github.com/yasharya2901/smart_divide/models"
	"gorm.io/gorm"
)

type ExpenseController struct {
	db *gorm.DB
}

func NewExpenseController(db *gorm.DB) *ExpenseController {
	return &ExpenseController{db: db}
}

func (ec *ExpenseController) CreateExpense(name string, totalAmount float64, eventID uint, paidByUserId uint) (*models.Expense, error) {
	// Create an expense
	expense := models.Expense{
		Name:        name,
		TotalAmount: totalAmount,
		EventID:     eventID,
		PaidByID:    paidByUserId,
	}
	if err := ec.db.Create(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (ec *ExpenseController) GetExpenses() ([]models.Expense, error) {
	// Get all expenses
	var expenses []models.Expense
	if err := ec.db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (ec *ExpenseController) GetExpenseByID(id uint) (*models.Expense, error) {
	// Get an expense by ID
	var expense models.Expense
	if err := ec.db.First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (ec *ExpenseController) UpdateExpense(id uint, name string, totalAmount float64, paidById uint) (*models.Expense, error) {
	// Update an expense
	var expense models.Expense

	if err := ec.db.First(&expense, id).Error; err != nil {
		return nil, err
	}

	expense.Name = name
	expense.TotalAmount = totalAmount
	expense.PaidByID = paidById

	if err := ec.db.Save(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil

}

func (ec *ExpenseController) AddExpensePerson(expenseId, personId uint, paidAmount float64, owedAmount float64) (*models.ExpensePerson, error) {
	// Add a person to an expense
	expensePerson := models.ExpensePerson{
		ExpenseID:  expenseId,
		PersonID:   personId,
		PaidAmount: paidAmount,
		OwedAmount: owedAmount,
	}
	if err := ec.db.Create(&expensePerson).Error; err != nil {
		return nil, err
	}
	return &expensePerson, nil
}

func (ec *ExpenseController) DeleteExpense(id uint) error {
	// Delete an expense
	var expense models.Expense
	if err := ec.db.First(&expense, id).Error; err != nil {
		return err
	}
	if err := ec.db.Delete(&expense).Error; err != nil {
		return err
	}
	return nil
}
