package services

import (
	"os"
	"strings"
)

var corsOrigins []string = []string{}

func SetOrigins() {
	originsStrings := os.Getenv("ORIGINS")
	if originsStrings == "" {
		corsOrigins = []string{"https://www.wilin.info", "https://wilin.info"}
		return
	}
	corsOrigins = strings.Split(originsStrings, ",")
}

func GetOrigins() []string {
	return corsOrigins
}
