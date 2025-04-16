package routes

import (
	"wilin/src/database"

	"github.com/gin-gonic/gin"
)

type Server struct {
	WordDao *database.WordDao
	UserDao *database.UserDao
}

type Header struct {
	Key, Value string
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

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
