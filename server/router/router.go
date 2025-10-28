package router

import (
	"context"
	"wilin/database/kalan"
	"wilin/database/users"
)

type ErrorJson struct {
	Error string `json:"error"`
}

func NewErrorJson(message string) ErrorJson {
	return ErrorJson{Error: message}
}

type Router struct {
	ctx          context.Context
	kalanQueries *kalan.Queries
	userQueries  *users.Queries
}

func New(
	ctx context.Context,
	kalanQueries *kalan.Queries,
	userQueries *users.Queries,
) *Router {
	return &Router{
		ctx:          ctx,
		kalanQueries: kalanQueries,
		userQueries:  userQueries,
	}
}
