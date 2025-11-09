package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"wilin.info/api/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func getDataSource() string {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("%s:%s@(%s)/%s?clientFoundRows=true&parseTime=true", username, password, address, dbName)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	db, err := sql.Open("mysql", getDataSource())
	if err != nil {
		log.Fatalf("Error opening database connection: %v\n", err)
	}

	server := server.New(db)
	server.Logger.Fatal(server.Start(":8080"))
}
