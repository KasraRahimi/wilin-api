package server

import (
	"context"
	"database/sql"
	"net/http"

	"wilin.info/api/database/kalan"
	"wilin.info/api/database/proposal"
	"wilin.info/api/database/recovery"
	"wilin.info/api/database/users"
	"wilin.info/api/server/router"
	"wilin.info/api/server/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const LOGGER_FORMAT = "\033[36m${time_custom}\033[0m | ${remote_ip} | \033[33m${method}\033[0m ${uri} | ${status} | ${latency_human}\n"
const TIME_FORMAT = "02-Jan-2006 15:04:05"

const MANUAL_LOGGER_FORMAT = "[${level}] | ${short_file}:${line} |${message}"

func newLoggerConfig(format string, timeFormat string) middleware.LoggerConfig {
	return middleware.LoggerConfig{
		Format:           format,
		CustomTimeFormat: timeFormat,
	}
}

func New(db *sql.DB) *echo.Echo {
	// initialize echo server
	server := echo.New()
	server.Logger.SetHeader(MANUAL_LOGGER_FORMAT)
	server.Logger.SetLevel(log.INFO)

	server.Use(
		middleware.LoggerWithConfig(
			newLoggerConfig(LOGGER_FORMAT, TIME_FORMAT),
		),
	)
	server.Use(middleware.Recover())

	// initialize router
	kalanQueries := kalan.New(db)
	usersQueries := users.New(db)
	proposalQueries := proposal.New(db)
	recoveryQueries := recovery.New(db)
	router := router.New(
		context.Background(),
		kalanQueries,
		usersQueries,
		proposalQueries,
		recoveryQueries,
	)

	// add preroute middleware
	services.SetOrigins()
	corsConfig := middleware.CORSConfig{
		AllowOrigins: services.GetOrigins(),
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderCacheControl,
			echo.HeaderXRequestedWith,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowCredentials: true,
	}
	server.Use(middleware.CORSWithConfig(corsConfig))
	server.Use(router.ExtractUserID)

	// add routes
	server.GET("/", router.HelloWorld)

	server.GET(
		"/kalan",
		router.GetAllKalan,
		router.VerifyPermissionsAll(services.PERMISSION_VIEW_WORD),
	)
	server.GET(
		"/kalan/paginated",
		router.GetKalanBySearch,
		router.VerifyPermissionsAll(services.PERMISSION_VIEW_WORD),
	)
	server.GET(
		"/kalan/:id",
		router.GetKalanByID,
		router.VerifyPermissionsAll(services.PERMISSION_VIEW_WORD),
	)

	server.POST(
		"/kalan",
		router.AddKalan,
		router.VerifyPermissionsAll(services.PERMISSION_ADD_WORD),
	)
	server.PUT(
		"/kalan",
		router.UpdateKalan,
		router.VerifyPermissionsAll(services.PERMISSION_MODIFY_WORD),
	)
	server.DELETE(
		"/kalan/:id",
		router.DeleteKalan,
		router.VerifyPermissionsAll(services.PERMISSION_DELETE_WORD),
	)

	server.GET(
		"/proposal",
		router.GetAllProposals,
		router.VerifyPermissionsAll(services.PERMISSION_VIEW_ALL_PROPOSAL),
	)
	server.GET(
		"/proposal/:id",
		router.GetProposalByID,
		router.VerifyPermissionsAny(services.PERMISSION_VIEW_SELF_PROPOSAL, services.PERMISSION_VIEW_ALL_PROPOSAL),
	)
	server.POST(
		"/proposal",
		router.PostProposal,
		router.VerifyPermissionsAll(services.PERMISSION_ADD_PROPOSAL),
	)
	server.PUT(
		"/proposal",
		router.UpdateProposal,
		router.VerifyPermissionsAny(services.PERMISSION_MODIFY_SELF_PROPOSAL, services.PERMISSION_MODIFY_ALL_PROPOSAL),
	)
	server.DELETE(
		"/proposal/:id",
		router.DeleteProposal,
		router.VerifyPermissionsAny(services.PERMISSION_DELETE_SELF_PROPOSAL, services.PERMISSION_DELETE_ALL_PROPOSAL),
	)
	server.GET(
		"/proposal/me",
		router.GetMyProposals,
		router.VerifyPermissionsAll(services.PERMISSION_VIEW_SELF_PROPOSAL),
	)

	server.POST("/signup", router.HandleSignUp)
	server.POST("/login", router.HandleLogin)
	server.GET("/me", router.GetMe)
	server.POST("/refresh", router.HandleRefresh)

	server.POST("/recovery", router.RequestRecovery)

	return server
}
