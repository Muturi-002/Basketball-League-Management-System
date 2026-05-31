package database

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ValidateUsername(username string, minLength, maxLength int) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < minLength || len(username) > maxLength {
		return fmt.Errorf("username must be between %d and %d characters", minLength, maxLength)
	}

	firstRune, _ := utf8.DecodeRuneInString(username)
	if !unicode.IsLetter(firstRune) {
		return fmt.Errorf("username must start with a letter")
	}

	for _, char := range username {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
			continue
		}
		return fmt.Errorf("username may only contain letters, numbers, and underscores")
	}

	return nil
}

func ValidatePassword(password string, minLength, maxLength int) error {
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < minLength || len(password) > maxLength {
		return fmt.Errorf("password must be between %d and %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateEmailAddress(emailAddress string) error {
	emailAddress = strings.TrimSpace(emailAddress)
	if emailAddress == "" {
		return fmt.Errorf("email address is required")
	}
	if len(emailAddress) < 6 || len(emailAddress) > 254 {
		return fmt.Errorf("email address must be between 6 and 254 characters")
	}
	if strings.ContainsAny(emailAddress, " \t\r\n") {
		return fmt.Errorf("email address cannot contain spaces")
	}
	parts := strings.Split(emailAddress, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !strings.Contains(parts[1], ".") {
		return fmt.Errorf("email address is invalid")
	}
	return nil
}

func ValidateLoginIdentifier(identifier string) error {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return fmt.Errorf("identifier is required")
	}
	if strings.Contains(identifier, "@") {
		return ValidateEmailAddress(identifier)
	}
	return ValidateUsername(identifier, 3, 32)
}