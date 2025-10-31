package router

import (
	"context"
	"errors"
	"strings"

	"wilin.info/api/database/kalan"
	"wilin.info/api/database/users"
)

type ErrorJson struct {
	Error string `json:"error"`
}

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrNoId          = errors.New("no id")
	ErrNoEntry       = errors.New("no entry")
	ErrNoPos         = errors.New("no pos")
	ErrNoGloss       = errors.New("no gloss")
	ErrNoUserID      = errors.New("no user id")

	ErrNoUserFromCtx = errors.New("user could not be fetched")
)

func NewErrorJson(message string) ErrorJson {
	return ErrorJson{Error: message}
}

// splitQuery takes a string and returns a slice
// with the string split by commas (,).
// If the string is empty, it will return an
// empty slice
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
