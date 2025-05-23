package database

import (
	"database/sql"
	"errors"
	"fmt"
	"wilin/src/database/roles"
)

type UserModel struct {
	Id           int
	Email        string
	Username     string
	PasswordHash string
	Role         roles.Role
}

const (
	idColumn       = "id"
	emailColumn    = "email"
	usernameColumn = "username"
)

var (
	ErrInvalidColumn = errors.New("invalid column")
)

type UserDao struct {
	Db *sql.DB
}

func (dao *UserDao) CreateUser(user *UserModel) (int64, error) {
	result, err := dao.Db.Exec(`
		INSERT INTO users (email, username, passwordHash, role) 
		VALUES (?, ?, ?, ?)
	`, user.Email, user.Username, user.PasswordHash, user.Role)

	if err != nil {
		return 0, fmt.Errorf("CreateUser, could not insert row: %w", err)
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("CreateUser, could not get last inserted row id: %w", err)
	}

	return lastId, nil
}

func (dao *UserDao) readUser(user *UserModel, column string) (*UserModel, error) {
	query := fmt.Sprintf("SELECT id, email, username, passwordHash, role FROM users WHERE %s = ?", column)
	var readUser UserModel
	var row *sql.Row

	switch column {
	case idColumn:
		row = dao.Db.QueryRow(query, user.Id)
	case emailColumn:
		row = dao.Db.QueryRow(query, user.Email)
	case usernameColumn:
		row = dao.Db.QueryRow(query, user.Username)
	default:
		return nil, ErrInvalidColumn
	}

	err := row.Scan(&readUser.Id, &readUser.Email, &readUser.Username, &readUser.PasswordHash, &readUser.Role)
	if err != nil {
		return nil, fmt.Errorf("readUser, could not scan row: %w", err)
	}

	return &readUser, nil
}

func (dao *UserDao) ReadUserById(id int) (*UserModel, error) {
	user, err := dao.readUser(&UserModel{Id: id}, idColumn)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) ReadUserByEmail(email string) (*UserModel, error) {
	user, err := dao.readUser(&UserModel{Email: email}, emailColumn)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) ReadUserByUsername(username string) (*UserModel, error) {
	user, err := dao.readUser(&UserModel{Username: username}, usernameColumn)
	if err != nil {
		return nil, err
	}
	return user, nil
}
