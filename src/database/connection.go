package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbname  = "wilin"
	address = "localhost:3306"
	driver  = "mysql"
)

func GetConnection() (*sql.DB, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	dataSource := fmt.Sprintf("%s:%s@(%s)/%s", username, password, address, dbname)
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, fmt.Errorf("GetConnection, failed at getting database connection: %w", err)
	}
	return db, nil
}

func CreateDictionaryTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id 		int PRIMARY KEY AUTO_INCREMENT,
			entry 	varchar(255) NOT NULL,
			pos 	varchar(255) NOT NULL,
			gloss	varchar(255) NOT NULL,
			notes 	varchar(2047) NOT NULL
		)
	`)

	if err != nil {
		return fmt.Errorf("CreateDictionaryTable, could not create dictionary table: %w", err)
	}
	return nil
}

func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id int PRIMARY KEY AUTO_INCREMENT,
		    email varchar(127) UNIQUE NOT NULL,
		    username varchar(31) UNIQUE NOT NULL,
		    passwordHash varchar(255) NOT NULL,
		    role varchar(255) NOT NULL
		)
	`)

	if err != nil {
		return fmt.Errorf("CreateUserTable, could not create dictionary table: %w", err)
	}
	return nil
}
