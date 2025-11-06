package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"wilin.info/api/database/users"
	"wilin.info/api/server/services"
)

func extractUserRole(
	ctx context.Context,
	userQueries *users.Queries,
	userIDInterface interface{},
) services.Role {
	userID, ok := userIDInterface.(int)
	if !ok {
		return services.ROLE_GUEST
	}

	user, err := userQueries.ReadUserByID(ctx, int32(userID))
	if err != nil {
		return services.ROLE_GUEST
	}

	return services.NewRole(user.Role)
}

func handleUnauthorized(ctx echo.Context, role services.Role) error {
	if role == services.ROLE_GUEST {
		return ctx.NoContent(http.StatusUnauthorized)
	} else {
		return ctx.NoContent(http.StatusForbidden)
	}
}

func (r *Router) VerifyPermissionsAll(perms ...services.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userIDInterface := ctx.Get("userID")
			role := extractUserRole(r.ctx, r.userQueries, userIDInterface)

			if role.CanAll(perms...) {
				return next(ctx)
			}

			return handleUnauthorized(ctx, role)
		}
	}
}

func (r *Router) VerifyPermissionsAny(perms ...services.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userIDInterface := ctx.Get("userID")
			role := extractUserRole(r.ctx, r.userQueries, userIDInterface)

			if role.CanAny(perms...) {
				return next(ctx)
			}

			return handleUnauthorized(ctx, role)
		}
	}
}
