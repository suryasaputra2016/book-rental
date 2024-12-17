package utils

import (
	"errors"
	"regexp"
	"unicode"
)

// IsEmailStringValid returns boolean of whether the string is in email format
func IsEmailStringValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// IsPasswordGood returns booleans of good password conditions
func IsPasswordGood(password string) error {
	var containNumber, containUpperCase, containPunctuation, eightOrMore bool
	index := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			containNumber = true
		case unicode.IsUpper(c):
			containUpperCase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			containPunctuation = true
		}
		index++
	}
	eightOrMore = index >= 8

	switch {
	case !containNumber:
		return errors.New("password must contain number")
	case !containUpperCase:
		return errors.New("password must contain upper case")
	case !containPunctuation:
		return errors.New("password must contain punctuation or special symbol")
	case !eightOrMore:
		return errors.New("password must contain eight or more characters")
	default:
		return nil
	}
}
