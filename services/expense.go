package services

import (
	"errors"

	"github.com/yasharya2901/smart_divide/models"
	"gorm.io/gorm"
)

type ExpenseService struct {
	db *gorm.DB
}

func NewExpenseService(db *gorm.DB) *ExpenseService {
	return &ExpenseService{db: db}
}

func (ec *ExpenseService) CreateExpense(name string, totalAmount float64, eventID, paidByUserId uint) (*models.Expense, error) {
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

func (ec *ExpenseService) GetExpenses() ([]models.Expense, error) {
	// Get all expenses
	var expenses []models.Expense
	if err := ec.db.Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (ec *ExpenseService) GetExpenseByID(id uint) (*models.Expense, error) {
	// Get an expense by ID
	var expense models.Expense
	if err := ec.db.First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (ec *ExpenseService) UpdateExpense(id uint, name string, totalAmount float64, paidById uint) (*models.Expense, error) {
	// Update an expense
	var expense models.Expense

	if err := ec.db.First(&expense, id).Error; err != nil {
		return nil, err
	}

	if name != "" {
		expense.Name = name
	}

	if totalAmount != 0 {
		expense.TotalAmount = totalAmount
	}

	if paidById != 0 {
		expense.PaidByID = paidById
	}

	if err := ec.db.Save(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil

}

func (ec *ExpenseService) AddExpensePerson(expenseId, personId uint) (*models.ExpensePerson, error) {
	// Add a person to an expense
	expensePerson := models.ExpensePerson{
		ExpenseID: expenseId,
		PersonID:  personId,
	}
	if err := ec.db.Create(&expensePerson).Error; err != nil {
		return nil, err
	}
	return &expensePerson, nil
}

func (ec *ExpenseService) GetExpensePeople(expenseId uint) ([]models.ExpensePerson, error) {
	// Get all people for an expense
	var expensePeople []models.ExpensePerson
	if err := ec.db.Where("expense_id = ?", expenseId).Find(&expensePeople).Error; err != nil {
		return nil, err
	}
	return expensePeople, nil
}

func (ec *ExpenseService) UpdateExpensePerson(expenseId, personId uint, paidAmount float64, owedAmount float64) (*models.ExpensePerson, error) {
	// Update an expense person
	var expensePerson models.ExpensePerson
	if err := ec.db.Where("expense_id = ? AND person_id = ?", expenseId, personId).First(&expensePerson).Error; err != nil {
		return nil, err
	}

	if paidAmount != 0 {
		expensePerson.PaidAmount = paidAmount
	}

	if owedAmount != 0 {
		expensePerson.OwedAmount = owedAmount
	}

	if err := ec.db.Save(&expensePerson).Error; err != nil {
		return nil, err
	}

	return &expensePerson, nil
}

func (ec *ExpenseService) DeleteExpensePerson(expenseId, personId uint) error {
	// Delete an expense person
	var expensePerson models.ExpensePerson
	if err := ec.db.Where("expense_id = ? AND person_id = ?", expenseId, personId).First(&expensePerson).Error; err != nil {
		return err
	}

	if err := ec.db.Delete(&expensePerson).Error; err != nil {
		return err
	}

	return nil
}

func (ec *ExpenseService) DeleteExpense(id uint) error {
	// Delete an expense
	var expense models.Expense
	if err := ec.db.First(&expense, id).Error; err != nil {
		return err
	}
	if err := ec.db.Delete(&expense).Error; err != nil {
		return err
	}

	// Also delete all splits for the expense
	if err := ec.db.Where("expense_id = ?", id).Delete(&models.ExpensePerson{}).Error; err != nil {
		return err
	}

	return nil
}

func (ec *ExpenseService) CheckExpenseConsistency(expenseId uint) error {
	// Check if the total amount of an expense is equal to the sum of the owed amounts
	var expense models.Expense
	if err := ec.db.First(&expense, expenseId).Error; err != nil {
		return err
	}

	var totalOwedAmount float64
	if err := ec.db.Model(&models.ExpensePerson{}).Where("expense_id = ?", expenseId).Select("sum(owed_amount)").Row().Scan(&totalOwedAmount); err != nil {
		return err
	}

	if totalOwedAmount != expense.TotalAmount {
		return errors.New("total owed amount does not match total amount")
	}

	return nil
}
