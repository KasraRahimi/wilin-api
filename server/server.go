package server

import (
	"context"
	"database/sql"

	"wilin.info/api/database/kalan"
	"wilin.info/api/database/users"
	"wilin.info/api/server/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const LOGGER_FORMAT = "\033[36m${time_custom}\033[0m | ${remote_ip} | \033[33m${method}\033[0m ${uri} | ${status} | ${latency_human}\n"
const TIME_FORMAT = "02-Jan-2006 15:04:05"

func newLoggerConfig(format string, timeFormat string) middleware.LoggerConfig {
	return middleware.LoggerConfig{
		Format:           format,
		CustomTimeFormat: timeFormat,
	}
}

func New(db *sql.DB) *echo.Echo {
	// initialize echo server
	server := echo.New()
	server.Use(
		middleware.LoggerWithConfig(
			newLoggerConfig(LOGGER_FORMAT, TIME_FORMAT),
		),
	)
	server.Use(middleware.Recover())

	// initialize router
	kalanQueries := kalan.New(db)
	usersQueries := users.New(db)
	router := router.New(context.Background(), kalanQueries, usersQueries)

	// add routes
	server.GET("/", router.HelloWorld)

	server.GET("/kalan", router.GetAllKalan)
	server.GET("/kalan/paginated", router.GetKalanBySearch)
	server.GET("/kalan/:id", router.GetKalanByID)
	server.POST("/kalan", router.AddKalan)
	server.PUT("/kalan", router.UpdateKalan)
	server.DELETE("/kalan/:id", router.DeleteKalan)

	server.POST("/signup", router.HandleSignUp)

	return server
}
