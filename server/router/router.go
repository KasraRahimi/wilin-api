package router

import (
	"context"
	"strings"
	"wilin/database/kalan"
	"wilin/database/users"
)

type ErrorJson struct {
	Error string `json:"error"`
}

func NewErrorJson(message string) ErrorJson {
	return ErrorJson{Error: message}
}

func splitQuery(query string) []string {
	if len(query) < 1 {
		return []string{}
	}

	queries := strings.Split(query, ",")
	for i := range queries {
		queries[i] = strings.TrimSpace(queries[i])
	}
	return queries
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
