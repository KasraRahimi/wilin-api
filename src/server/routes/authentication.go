package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"wilin/src/database"
	"wilin/src/server/utils"
)

type LoginFields struct {
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

func (dto *UserDTO) FromUserModel(userModel *database.UserModel) *UserDTO {
	return &UserDTO{
		Id:       userModel.Id,
		Email:    userModel.Email,
		Username: userModel.Username,
		Password: "", // We don't want to send the password hash back to the frontend
		Role:     userModel.Role,
	}
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

	time.Sleep(1000 * time.Millisecond)

	if loginFields.Username == "kawa" || loginFields.Password == "wikimokalan" {
		ctx.JSON(http.StatusOK, gin.H{"token": "nimilen"})
		return
	}

	ctx.String(http.StatusUnauthorized, "Invalid login information")
}

func (s *Server) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := "nimilen"
		authHeader := ctx.GetHeader("Authorization")
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) != 2 {
			ctx.String(http.StatusUnauthorized, "Incorrect token")
			ctx.Abort()
			return
		}
		if authHeaders[1] != token {
			ctx.String(http.StatusUnauthorized, "Incorrect token")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
