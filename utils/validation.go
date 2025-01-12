package utils

import (
	"regexp"
	"strings"
)

func ValidatePhoneNumber(phone string) bool {
	// Check if empty
	if len(phone) == 0 {
		return false
	}

	// Must start with '+'
	if !strings.HasPrefix(phone, "+") {
		return false
	}

	// Minimum length check ('+' and at least 7 digits)
	if len(phone) < 8 {
		return false
	}

	// Check if rest are digits using regex
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
	return phoneRegex.MatchString(phone)
}

func ValidateEmail(email string) bool {
	// Check if empty
	if len(email) == 0 {
		return false
	}

	// Check length (typical max email length is 254 characters)
	if len(email) > 254 {
		return false
	}

	// Email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
