package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model           // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string    `gorm:"type:varchar(255);not null"` // Event name
	People     []Person  `gorm:"many2many:event_people"`     // Many-to-many relationship with People
	Expenses   []Expense `gorm:"foreignKey:EventID"`         // One-to-many relationship with Expense
}

type Expense struct {
	gorm.Model                  // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name        string          `gorm:"type:varchar(255);not null"`  // Expense name
	TotalAmount float64         `gorm:"type:decimal(10,2);not null"` // Total expense amount
	EventID     uint            `gorm:"not null"`                    // Foreign key to Event
	PaidByID    uint            `gorm:"not null"`                    // Foreign key to People
	PaidBy      Person          `gorm:"foreignKey:PaidByID"`         // Reference to the person who paid
	Splits      []ExpensePerson `gorm:"foreignKey:ExpenseID"`        // Splits for the expense
}

type Person struct {
	gorm.Model                             // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name                   string          `gorm:"type:varchar(255);not null"` // Person name
	Contact                string          `gorm:"type:varchar(50)"`           // Contact number
	Email                  string          `gorm:"type:varchar(255);unique"`   // Email (unique constraint)
	Password               string          `gorm:"type:varchar(255);not null"` // Hashed password
	RefreshToken           string          `gorm:"type:text"`                  // Refresh token
	RefreshTokenExpiryDate *time.Time      `gorm:"type:timestamp"`             // Refresh token expiry date
	Events                 []Event         `gorm:"many2many:event_people"`     // Many-to-many relationship with Event
	Expenses               []ExpensePerson `gorm:"foreignKey:PersonID"`        // Splits for expenses
}

type ExpensePerson struct {
	gorm.Model         // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	ExpenseID  uint    `gorm:"not null"`             // Foreign key to Expense
	PersonID   uint    `gorm:"not null"`             // Foreign key to Person
	PaidAmount float64 `gorm:"type:decimal(10,2)"`   // Amount paid by the person
	OwedAmount float64 `gorm:"type:decimal(10,2)"`   // Amount owed by the person
	Expense    Expense `gorm:"foreignKey:ExpenseID"` // Reference to the expense
	Person     Person  `gorm:"foreignKey:PersonID"`  // Reference to the person
}
