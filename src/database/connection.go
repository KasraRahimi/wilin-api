package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "kawa"
	password = "password"
	dbname   = "wilin"
	address  = "localhost:3306"
	driver   = "mysql"
)

func GetConnection() (*sql.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@(%s)/%s", username, password, address, dbname)
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, fmt.Errorf("GetConnection, failed at getting database connection: %w", err)
	}
	return db, nil
}

func CreateDictionaryTable() error {
	db, err := GetConnection()
	if err != nil {
		return fmt.Errorf("CreateDictionaryTable, could not connect to database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`
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
