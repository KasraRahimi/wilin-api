package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Router) HelloWorld(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, World!")
}
