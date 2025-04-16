package server

import (
	"wilin/src/database"
	"wilin/src/server/routes"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	server := routes.Server{
		WordDao: &database.WordDao{},
		UserDao: &database.UserDao{},
	}
	router := gin.Default()
	router.Use(server.CorsMiddleware())

	router.GET("/kalan", server.HandleGetKalan)
	router.GET("/kalan/:id", server.HandleGetKalanById)
	router.POST("/login", server.HandleLogin)
	router.POST("/signup", server.HandleSignup)

	restricted := router.Group("")
	restricted.Use(server.Authentication())

	restricted.POST("/kalan", server.HandlePostKalan)
	restricted.DELETE("/kalan/:id", server.HandleDeleteKalan)
	restricted.PUT("/kalan", server.HandlePutKalan)

	return router
}
