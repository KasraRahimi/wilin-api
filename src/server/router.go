package server

import (
	"wilin/src/database"
	"wilin/src/server/routes"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	wordDao := database.WordDao{}
	server := routes.Server{
		WordDao: &wordDao,
	}
	router := gin.Default()
	router.GET("/kalan", server.HandleGetKalan)
	router.GET("/kalan/:id", server.HandleGetKalanById)
	router.POST("/kalan", server.HandlePostKalan)
	return router
}
