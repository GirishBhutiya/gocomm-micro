package util

import (
	"log"
	"regexp"
)

func CheckPasswordValidity(password string) bool {
	if regexp.MustCompile(`\s`).MatchString(password) {
		log.Println("in regex")
		return false
	}
	if len(password) < 8 {
		log.Println("in length")
		return false

	}
	return true
}
