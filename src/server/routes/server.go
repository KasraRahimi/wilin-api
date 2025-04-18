package routes

import (
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
