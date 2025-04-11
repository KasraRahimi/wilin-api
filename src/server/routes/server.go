package routes

import (
	"net/http"
	"strings"
	"wilin/src/database"

	"github.com/gin-gonic/gin"
)

type Server struct {
	WordDao *database.WordDao
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
