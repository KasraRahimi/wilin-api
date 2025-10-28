package server

import (
	"context"
	"database/sql"
	"net/http"
	"wilin/database/kalan"
	"wilin/database/users"
	"wilin/server/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func helloWorld(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, World!")
}

const LOGGER_FORMAT = "${time_custom} | ${remote_ip} | ${method} ${path} | ${status} | ${latency_human}\n"
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

	return server
}
