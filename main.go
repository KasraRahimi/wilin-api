package main

import (
	"github.com/joho/godotenv"
	"log"
	"wilin/src/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	router := server.New()
	err = router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
