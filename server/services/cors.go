package services

import (
	"os"
	"strings"
)

var DEFAULT_ORIGINS = []string{"https://www.wilin.info", "https://wilin.info"}
var corsOrigins = DEFAULT_ORIGINS

func SetOrigins() {
	originsStrings := os.Getenv("ORIGINS")
	if originsStrings == "" {
		corsOrigins = DEFAULT_ORIGINS
		return
	}

	parts := strings.Split(originsStrings, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	corsOrigins = parts
}

func GetOrigins() []string {
	return corsOrigins
}
