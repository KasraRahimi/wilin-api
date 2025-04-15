package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type LoginFields struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
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
