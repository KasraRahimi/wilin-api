package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"wilin/src/database"
	"wilin/src/database/permissions"
	"wilin/src/database/roles"
	"wilin/src/server/utils"
)

const TIME_TO_EXPIRE_MINUTES = 60 * 12 // 12 hours before a token expires

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
	Id       int    `json:"id" form:"id"`
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

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
			_ = utils.IsPasswordAndHashSame(loginFields.Password, loginFields.Password)
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

	token, err := utils.GenerateToken(strconv.Itoa(user.Id), TIME_TO_EXPIRE_MINUTES)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto := UserDTO{}
	dto.FromUserModel(user)
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": dto})
}

func (s *Server) HandleSignup(ctx *gin.Context) {
	var signUpFields SignUpFields
	if err := ctx.ShouldBind(&signUpFields); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusBadRequest, "Incorrectly formatted")
		return
	}

	userDTo := UserDTO{
		Id:       0,
		Email:    signUpFields.Email,
		Username: signUpFields.Username,
		Password: signUpFields.Password,
		Role:     roles.USER,
	}

	if !utils.IsValidEmail(userDTo.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	userModel, err := userDTo.ToUserModel()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = s.UserDao.ReadUserByEmail(userModel.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already taken"})
		return
	}

	_, err = s.UserDao.ReadUserByUsername(userModel.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
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

func (s *Server) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) != 2 {
			ctx.Next()
			return
		}
		tokenString := authHeaders[1]
		token, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}
		ctx.Set("uid", token.Id)
		ctx.Next()
	}
}

func (s *Server) VerifyPermissions(perms ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var role string
		uid, exists := ctx.Get("uid")
		if exists {
			id, err := strconv.Atoi(uid.(string))
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}
			user, err := s.UserDao.ReadUserById(id)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
				ctx.Abort()
				return
			}
			role = user.Role
		} else {
			role = roles.NON_USER
		}
		for _, permission := range perms {
			if !permissions.CanRolePermission(role, permission) {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
				ctx.Abort()
			}
		}
		ctx.Next()
	}
}
