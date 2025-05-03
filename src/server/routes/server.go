package routes

import (
	"net/http"
	"wilin/src/database"
	"wilin/src/server/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	WordDao          *database.WordDao
	UserDao          *database.UserDao
	fakePasswordHash string
}

func (s *Server) GenerateFakePassword() {
	hash, err := utils.GeneratePasswordHash("password")
	if err != nil {
		s.fakePasswordHash = "$2a$13$jidvuZmt0Bji5yuIckM4jOFKwF062Lt3M8M.A4uGSKIkS3r8BA00O"
		return
	}
	s.fakePasswordHash = hash
}

type Header struct {
	Key, Value string
}

func (s *Server) CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		originHeader := Header{Key: "Access-Control-Allow-Origin", Value: "https://www.wilin.info"}
		if origin := ctx.Request.Header.Get("Origin"); origin == "http://localhost:3000" {
			originHeader.Value = origin
		}

		headers := [...]Header{
			originHeader,
			{"Access-Control-Allow-Credentials", "true"},
			{"Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Cache-Control, X-Custom-Header"},
			{"Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE"},
		}

		for _, header := range headers {
			ctx.Writer.Header().Set(header.Key, header.Value)
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}
