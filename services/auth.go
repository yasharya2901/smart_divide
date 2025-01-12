package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/yasharya2901/smart_divide/models"
	"github.com/yasharya2901/smart_divide/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db            *gorm.DB
	peopleService *PeopleService
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db, peopleService: NewPeopleService(db)}
}

func (as *AuthService) Login(email, password string) (string, string, error) {
	// Login a user
	var person models.Person
	result := as.peopleService.db.Where("email = ?", email).First(&person)
	if result.Error != nil {
		return "", "", errors.New("invalid email")
	}

	if valid, err := utils.ComparePasswords(person.Password, password); err != nil || !valid {
		return "", "", errors.New("invalid password")
	}

	// Expiry time for the refresh and access token
	refreshTokenExpiry, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	accessTokenExpiry, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	// Generate an access token
	accessToken, _, err := utils.GenerateToken(person.ID, person.Email, time.Minute*time.Duration(accessTokenExpiry), os.Getenv("JWT_ACCESS_SECRET"))
	if err != nil {
		return "", "", err
	}

	// If the user already has a refresh token, return it
	if person.RefreshToken != "" {
		return accessToken, person.RefreshToken, nil
	}

	refreshToken, expiryTime, err := utils.GenerateToken(person.ID, person.Email, time.Hour*24*time.Duration(refreshTokenExpiry), os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil {
		return "", "", err
	}

	as.peopleService.UpdatePerson(person.ID, "", "", "", refreshToken, expiryTime)
	return accessToken, refreshToken, nil

}

func (as *AuthService) Register(name, phoneNumber, email, password string) (string, string, error) {
	// Register a user

	// Check if the email is already registered
	var person models.Person
	result := as.peopleService.db.Where("email = ?", email).First(&person)
	if result.Error == nil {
		return "", "", errors.New("email already registered")
	}

	if result.Error != gorm.ErrRecordNotFound {
		return "", "", result.Error
	}

	// Check if the phone number is already registered
	result = as.peopleService.db.Where("contact = ?", phoneNumber).First(&person)
	if result.Error == nil {
		return "", "", errors.New("phone number already registered")
	}

	if result.Error != gorm.ErrRecordNotFound {
		return "", "", result.Error
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", "", err
	}

	// Create a person
	person = models.Person{
		Name:     name,
		Contact:  phoneNumber,
		Email:    email,
		Password: hashedPassword,
	}

	if err := as.peopleService.db.Create(&person).Error; err != nil {
		return "", "", err
	}

	// Login the user
	accessToken, refreshToken, err := as.Login(email, password)
	return accessToken, refreshToken, err
}

func (as *AuthService) ChangePassword(email, oldPassword, newPassword string) error {
	// Change password of a person
	var person models.Person

	hashedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := as.peopleService.db.Where("email = ?", email).First(&person).Error; err != nil {
		return err
	}

	passwordMatched, err := utils.ComparePasswords(person.Password, oldPassword)
	if err != nil {
		return err
	}

	if !passwordMatched {
		return errors.New("incorrect password")
	}

	person.Password = hashedNewPassword
	if err := as.peopleService.db.Save(&person).Error; err != nil {
		return err
	}

	return nil
}

func (as *AuthService) RegenerateAccessToken(refreshToken string) (string, error) {
	// Regenerate an access token
	claims, err := utils.ValidateToken(refreshToken, os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil {
		return "", err
	}

	// If the refresh token is valid, check if the user exists
	_, err = as.peopleService.GetPersonByID(claims.UserID)
	if err != nil {
		return "", err
	}

	accessTokenExpiry, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return "", err
	}

	accessToken, _, err := utils.GenerateToken(claims.UserID, claims.Email, time.Minute*time.Duration(accessTokenExpiry), os.Getenv("JWT_ACCESS_SECRET"))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (as *AuthService) Logout(refreshToken string) error {
	// Logout a user
	claims, err := utils.ValidateToken(refreshToken, os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil {
		return err
	}

	// If the refresh token is valid, check if the user exists
	person, err := as.peopleService.GetPersonByID(claims.UserID)
	if err != nil {
		return err
	}

	person.RefreshToken = ""
	person.RefreshTokenExpiryDate = nil
	if err := as.peopleService.db.Save(&person).Error; err != nil {
		return err
	}

	return nil
}
