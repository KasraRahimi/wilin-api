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

type Header struct {
	Key, Value string
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

func (s *Server) CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headers := [...]Header{
			{"Access-Control-Allow-Origin", "*"},
			{"Access-Control-Allow-Credentials", "true"},
			{"Access-Control-Allow-Headers", "*"},
			{"Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE"},
		}
		for _, header := range headers {
			ctx.Writer.Header().Set(header.Key, header.Value)
		}
		ctx.Next()
	}
}
