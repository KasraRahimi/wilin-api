package server

import (
	"wilin/src/database"
	"wilin/src/database/permissions"
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
	router.Use(server.Authentication())
	router.POST("/login", server.HandleLogin)
	router.POST("/signup", server.HandleSignup)

	router.GET("/kalan", server.VerifyPermissions(permissions.VIEW_WORD), server.HandleGetKalan)
	router.GET("/kalan/:id", server.VerifyPermissions(permissions.VIEW_WORD), server.HandleGetKalanById)

	router.POST("/kalan", server.VerifyPermissions(permissions.ADD_WORD), server.HandlePostKalan)
	router.DELETE("/kalan/:id", server.VerifyPermissions(permissions.DELETE_WORD), server.HandleDeleteKalan)
	router.PUT("/kalan", server.VerifyPermissions(permissions.MODIFY_WORD), server.HandlePutKalan)

	return router
}
