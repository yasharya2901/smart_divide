package services

import (
	"errors"

	"github.com/yasharya2901/smart_divide/models"
	"github.com/yasharya2901/smart_divide/utils"
	"gorm.io/gorm"
)

type PeopleService struct {
	db *gorm.DB
}

func NewPeopleService(db *gorm.DB) *PeopleService {
	return &PeopleService{db: db}
}

func (ps *PeopleService) CreatePerson(name, contact, email string) (*models.Person, error) {
	// Create a person
	person := models.Person{
		Name:    name,
		Contact: contact,
		Email:   email,
	}
	if err := ps.db.Create(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (ps *PeopleService) GetPeople() ([]models.Person, error) {
	// Get all people
	var people []models.Person
	if err := ps.db.Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (ps *PeopleService) GetPersonByID(id uint) (*models.Person, error) {
	// Get a person by ID
	var person models.Person
	if err := ps.db.First(&person, id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (ps *PeopleService) UpdatePerson(id uint, name, contact, email string) (*models.Person, error) {
	// Update a person
	var person models.Person

	if err := ps.db.First(&person, id).Error; err != nil {
		return nil, err
	}

	if name != "" {
		person.Name = name
	}
	if contact != "" {
		person.Contact = contact
	}
	if email != "" {
		person.Email = email
	}

	if err := ps.db.Save(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (ps *PeopleService) GetPeopleByContacts(contacts []string) ([]models.Person, error) {
	// Get people by contacts
	var people []models.Person
	if err := ps.db.Where("contact IN ?", contacts).Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (ps *PeopleService) GetPeopleByEmails(emails []string) ([]models.Person, error) {
	// Get people by emails
	var people []models.Person
	if err := ps.db.Where("email IN ?", emails).Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (ps *PeopleService) Authenticate(email, password string) (interface{}, error) {
	// Authenticate a person
	var person models.Person

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	if err := ps.db.Where("email = ?", email).First(&person).Error; err != nil {
		return nil, err
	}

	passwordMatched, err := utils.ComparePasswords(hashedPassword, person.Password)
	if err != nil {
		return nil, err
	}

	if !passwordMatched {
		return nil, errors.New("incorrect password")
	}

	return &struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Contact string `json:"contact"`
	}{
		Name:    person.Name,
		Email:   person.Email,
		Contact: person.Contact,
	}, nil
}

func (ps *PeopleService) ChangePassword(email, oldPassword, newPassword string) error {
	// Change password of a person
	var person models.Person

	hashedOldPassword, err := utils.HashPassword(oldPassword)
	if err != nil {
		return err
	}

	hashedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := ps.db.Where("email = ?", email).First(&person).Error; err != nil {
		return err
	}

	passwordMatched, err := utils.ComparePasswords(hashedOldPassword, person.Password)
	if err != nil {
		return err
	}

	if !passwordMatched {
		return errors.New("incorrect password")
	}

	person.Password = hashedNewPassword
	if err := ps.db.Save(&person).Error; err != nil {
		return err
	}

	return nil
}
