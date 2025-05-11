package server

import (
	"database/sql"
	"wilin/src/database"
	"wilin/src/database/permissions"
	"wilin/src/server/routes"

	"github.com/gin-gonic/gin"
)

func New(db *sql.DB) (*gin.Engine, error) {
	server := routes.Server{
		WordDao:     &database.WordDao{Db: db},
		UserDao:     &database.UserDao{Db: db},
		ProposalDao: &database.ProposalDao{Db: db},
	}
	server.GenerateFakePassword()

	router := gin.Default()
	router.Use(server.CorsMiddleware())
	router.Use(server.Authentication())

	router.GET("/me", server.HandleMe)

	router.POST("/login", server.HandleLogin)
	router.POST("/signup", server.HandleSignup)
	router.POST("/refresh", server.HandleRefresh)

	router.GET("/kalan", server.VerifyPermissionsAll(permissions.VIEW_WORD), server.HandleGetKalan)
	router.GET("/kalan/paginated", server.VerifyPermissionsAll(permissions.VIEW_WORD), server.HandleGetKalanPaginated)
	router.GET("/kalan/:id", server.VerifyPermissionsAll(permissions.VIEW_WORD), server.HandleGetKalanById)

	router.POST("/kalan", server.VerifyPermissionsAll(permissions.ADD_WORD), server.HandlePostKalan)
	router.DELETE("/kalan/:id", server.VerifyPermissionsAll(permissions.DELETE_WORD), server.HandleDeleteKalan)
	router.PUT("/kalan", server.VerifyPermissionsAll(permissions.MODIFY_WORD), server.HandlePutKalan)

	router.POST("/proposal", server.VerifyPermissionsAll(permissions.ADD_PROPOSAL), server.HandlePostProposal)
	router.GET("/proposal", server.VerifyPermissionsAll(permissions.VIEW_ALL_PROPOSAL), server.HandleGetAllProposals)
	router.GET("/proposal/me", server.VerifyPermissionsAll(permissions.VIEW_SELF_PROPOSAL), server.HandleGetMyProposals)
	router.DELETE(
		"/proposal/:id",
		server.VerifyPermissionsAny(permissions.DELETE_ALL_PROPOSAL, permissions.DELETE_SELF_PROPOSAL),
		server.HandleDeleteProposal,
	)

	return router, nil
}
