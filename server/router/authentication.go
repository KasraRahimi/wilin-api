package router

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"wilin.info/api/database/users"
	"wilin.info/api/server/services"
)

const TIME_TO_AUTH_EXPIRE_MINUTES = 15
const TIME_TO_REFRESH_EXPIRE_MINUTES = 60 * 24 * 31

type LoginFields struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type SignUpFields struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserDTO struct {
	ID       int    `json:"id" form:"id"`
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

func NewUserDTO(id int, email string, username string, password string, role string) UserDTO {
	return UserDTO{
		ID:       id,
		Email:    email,
		Username: username,
		Password: password,
		Role:     role,
	}
}

type RefreshTokenDTO struct {
	Token string `json:"token"`
}

const (
	MAX_EMAIL_LENGTH    = 120
	MAX_USERNAME_LENGTH = 30
	MAX_PASSWORD_LENGTH = 60
	MIN_PASSWORD_LENGTH = 8
)

const (
	InvalidForm = "invalid format"

	InvalidEmail = "invalid email"
	EmailTaken   = "email taken"
	EmailTooLong = "email too long"

	NoUsername       = "no username"
	UsernameHasSpace = "username has space"
	UsernameTaken    = "username taken"
	UsernameTooLong  = "username too long"

	PasswordTooShort = "password too short"
	PasswordTooLong  = "password too long"
)

func validateSignUpFields(ctx context.Context, userQueries *users.Queries, fields SignUpFields) (int, error) {
	if !services.IsValidEmail(fields.Email) {
		return http.StatusBadRequest, errors.New(InvalidEmail)
	}
	if len(fields.Email) > MAX_EMAIL_LENGTH {
		return http.StatusBadRequest, errors.New(EmailTooLong)
	}
	if len(fields.Username) == 0 {
		return http.StatusBadRequest, errors.New(NoUsername)
	}
	if len(fields.Username) > MAX_USERNAME_LENGTH {
		return http.StatusBadRequest, errors.New(UsernameTooLong)
	}
	if strings.Contains(fields.Username, " ") {
		return http.StatusBadRequest, errors.New(UsernameHasSpace)
	}
	if len(fields.Password) < MIN_PASSWORD_LENGTH {
		return http.StatusBadRequest, errors.New(PasswordTooShort)
	}
	if len(fields.Password) > MAX_PASSWORD_LENGTH {
		return http.StatusBadRequest, errors.New(PasswordTooLong)
	}

	_, err := userQueries.ReadUserByEmail(ctx, fields.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		return http.StatusConflict, errors.New(EmailTaken)
	}

	_, err = userQueries.ReadUserByUsername(ctx, fields.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		return http.StatusConflict, errors.New(UsernameTaken)
	}

	return http.StatusOK, nil
}

func (r *Router) HandleSignUp(ctx echo.Context) error {
	signUpFields := new(SignUpFields)

	err := ctx.Bind(signUpFields)
	if err != nil {
		errJSON := NewErrorJson(InvalidForm)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	statusCode, err := validateSignUpFields(r.ctx, r.userQueries, *signUpFields)
	if err != nil {
		errJSON := NewErrorJson(err.Error())
		return ctx.JSON(statusCode, errJSON)
	}

	passwordHash, err := services.GeneratePasswordHash(signUpFields.Password)
	if err != nil {
		errJSON := NewErrorJson("something went wrong")
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	params := users.CreateUserParams{
		Email:    signUpFields.Email,
		Username: signUpFields.Username,
		Password: passwordHash,
		Role:     services.RoleUser.String(),
	}
	result, err := r.userQueries.CreateUser(r.ctx, params)
	if err != nil {
		errJSON := NewErrorJson("failed to create user")
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		errJSON := NewErrorJson("something went wrong")
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userDTO := NewUserDTO(
		int(userID),
		signUpFields.Email,
		signUpFields.Username,
		"",
		services.RoleUser.String(),
	)

	return ctx.JSON(http.StatusCreated, userDTO)
}
