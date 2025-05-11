package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wilin/src/database"
	"wilin/src/database/permissions"
	"wilin/src/database/roles"
	"wilin/src/server/utils"

	"github.com/gin-gonic/gin"
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
	Id       int        `json:"id" form:"id"`
	Email    string     `json:"email" form:"email"`
	Username string     `json:"username" form:"username"`
	Password string     `json:"password" form:"password"`
	Role     roles.Role `json:"role" form:"role"`
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

func (dto *UserDTO) FromUserModel(userModel *database.UserModel) {
	dto.Id = userModel.Id
	dto.Email = userModel.Email
	dto.Username = userModel.Username
	dto.Role = userModel.Role
}

func (dto *UserDTO) ToUserModel() (*database.UserModel, error) {
	passwordHash, err := utils.GeneratePasswordHash(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("ToUserModel, generating password hash: %w", err)
	}
	return &database.UserModel{
		Id:           dto.Id,
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Role:         dto.Role,
	}, nil
}

func (s *Server) HandleLogin(ctx *gin.Context) {
	var loginFields LoginFields
	if err := ctx.ShouldBind(&loginFields); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusBadRequest, "Incorrectly formatted")
		return
	}

	user, err := s.UserDao.ReadUserByUsername(loginFields.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = utils.IsPasswordAndHashSame(loginFields.Password, s.fakePasswordHash)
			// still compute a hash so the user cannot tell if the username or password is incorrect
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !utils.IsPasswordAndHashSame(loginFields.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	authToken, err := utils.GenerateToken("auth", strconv.Itoa(user.Id), TIME_TO_AUTH_EXPIRE_MINUTES)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken("refresh", strconv.Itoa(user.Id), TIME_TO_REFRESH_EXPIRE_MINUTES)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto := UserDTO{}
	dto.FromUserModel(user)
	ctx.JSON(http.StatusOK, gin.H{"authToken": authToken, "refreshToken": refreshToken, "user": dto})
}

func (s *Server) validateSignUpFields(fields SignUpFields) (int, error) {
	if !utils.IsValidEmail(fields.Email) {
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

	_, err := s.UserDao.ReadUserByEmail(fields.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		return http.StatusConflict, errors.New(EmailTaken)
	}

	_, err = s.UserDao.ReadUserByUsername(fields.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		return http.StatusConflict, errors.New(UsernameTaken)
	}

	return http.StatusOK, nil
}

func (s *Server) HandleSignup(ctx *gin.Context) {
	var signUpFields SignUpFields
	if err := ctx.ShouldBind(&signUpFields); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": InvalidForm})
		return
	}

	statusCode, err := s.validateSignUpFields(signUpFields)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	userDTo := UserDTO{
		Id:       0,
		Email:    signUpFields.Email,
		Username: signUpFields.Username,
		Password: signUpFields.Password,
		Role:     roles.USER,
	}

	userModel, err := userDTo.ToUserModel()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId, err := s.UserDao.CreateUser(userModel)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDTo.Password = ""
	userDTo.Id = int(userId)

	ctx.JSON(http.StatusCreated, userDTo)
}

func (s *Server) HandleMe(ctx *gin.Context) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userDto := UserDTO{}
	userDto.FromUserModel(user)
	ctx.JSON(http.StatusOK, userDto)
}

func (s *Server) HandleRefresh(ctx *gin.Context) {
	var refreshToken RefreshTokenDTO
	if err := ctx.ShouldBind(&refreshToken); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": InvalidForm})
		return
	}
	token, err := utils.ParseToken(refreshToken.Token)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if utils.IsExpired(token.Exp) || token.TokenType != "refresh" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	newAuthToken, err := utils.GenerateToken("auth", token.Sub, TIME_TO_AUTH_EXPIRE_MINUTES)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"token": newAuthToken})
}

func (s *Server) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		tokenString := authHeaders[1]
		token, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		if utils.IsExpired(token.Exp) || token.TokenType != "auth" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("uid", token.Sub)
		ctx.Next()
	}
}

func (s *Server) VerifyPermissionsAll(perms ...permissions.Permission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := s.getUserFromContext(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		var role roles.Role
		if user == nil {
			role = roles.NON_USER
		} else {
			role = user.Role
		}

		for _, permission := range perms {
			if !permissions.CanRolePermission(role, permission) {
				if user == nil {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
					ctx.Abort()
					return
				}
				ctx.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

func (s *Server) VerifyPermissionsAny(perms ...permissions.Permission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := s.getUserFromContext(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		var role roles.Role
		if user == nil {
			role = roles.NON_USER
		} else {
			role = user.Role
		}

		for _, permission := range perms {
			if permissions.CanRolePermission(role, permission) {
				ctx.Next()
				return
			}
		}
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		ctx.Abort()
	}
}

func (s *Server) getUserFromContext(ctx *gin.Context) (*database.UserModel, error) {
	uid, exists := ctx.Get("uid")
	if !exists {
		return nil, nil
	}

	id, err := strconv.Atoi(uid.(string))
	if err != nil {
		return nil, err
	}

	user, err := s.UserDao.ReadUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
