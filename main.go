package main

import (
	"github.com/joho/godotenv"
	"log"
	"wilin/src/database"
	"wilin/src/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	db, err := database.GetConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	err = database.CreateDictionaryTable(db)
	if err != nil {
		log.Fatalf("Error creating dictionary table: %v\n", err)
	}

	err = database.CreateUserTable(db)
	if err != nil {
		log.Fatalf("Error creating user table: %v\n", err)
	}

	router, err := server.New(db)
	if err != nil {
		log.Fatalf("Error creating server: %v\n", err)
	}
	err = router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
