package router

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
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
	Password string `json:"password" form:"password,omitempty"`
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

type TokensDTO struct {
	AuthToken    string `json:"authToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type LoginReturnDTO struct {
	User UserDTO `json:"user"`
	TokensDTO
}

const (
	MAX_EMAIL_LENGTH    = 120
	MAX_USERNAME_LENGTH = 30
	MAX_PASSWORD_LENGTH = 60
	MIN_PASSWORD_LENGTH = 8
)

const (
	InvalidForm  = "invalid format"
	InvalidLogin = "invalid username or password"
	ServerError  = "something went wrong"

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

func (r *Router) ExtractUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			return next(ctx)
		}

		tokenString := authParts[1]
		claims, err := services.ParseToken(tokenString)
		if err != nil {
			return next(ctx)
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return next(ctx)
		}

		ctx.Set("userID", userID)
		return next(ctx)
	}
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
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	params := users.CreateUserParams{
		Email:    signUpFields.Email,
		Username: signUpFields.Username,
		Password: passwordHash,
		Role:     services.ROLE_USER.String(),
	}
	result, err := r.userQueries.CreateUser(r.ctx, params)
	if err != nil {
		errJSON := NewErrorJson("failed to create user")
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userDTO := NewUserDTO(
		int(userID),
		signUpFields.Email,
		signUpFields.Username,
		"",
		services.ROLE_USER.String(),
	)

	return ctx.JSON(http.StatusCreated, userDTO)
}

func (r *Router) HandleLogin(ctx echo.Context) error {
	loginFields := new(LoginFields)

	err := ctx.Bind(loginFields)
	if err != nil {
		errJSON := NewErrorJson(InvalidForm)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	user, err := r.userQueries.ReadUserByUsername(r.ctx, loginFields.Username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errJSON := NewErrorJson(ServerError)
			return ctx.JSON(http.StatusInternalServerError, errJSON)
		}

		// run a fake hash to simulate invalid password
		services.FakeHashCompare()

		errJSON := NewErrorJson(InvalidLogin)
		return ctx.JSON(http.StatusUnauthorized, errJSON)
	}

	isCorrectPassword := services.IsPasswordAndHashSame(loginFields.Password, user.Password)

	if !isCorrectPassword {
		errJSON := NewErrorJson(InvalidLogin)
		return ctx.JSON(http.StatusUnauthorized, errJSON)
	}

	userID := strconv.Itoa(int(user.ID))

	authToken, err := services.GenerateToken("authToken", userID, TIME_TO_AUTH_EXPIRE_MINUTES)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	refreshToken, err := services.GenerateToken("refreshToken", userID, TIME_TO_REFRESH_EXPIRE_MINUTES)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	tokensDTO := TokensDTO{AuthToken: authToken, RefreshToken: refreshToken}

	userDTO := NewUserDTO(
		int(user.ID),
		user.Email,
		user.Username,
		"",
		user.Role,
	)

	loginReturnDTO := LoginReturnDTO{User: userDTO, TokensDTO: tokensDTO}

	return ctx.JSON(http.StatusOK, loginReturnDTO)
}
