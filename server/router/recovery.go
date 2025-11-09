package router

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"wilin.info/api/database/recovery"
	"wilin.info/api/server/services"
)

type RequestRecoveryDTO struct {
	Email string `json:"email" form:"email"`
}

type ChangePasswordDTO struct {
	Password string `json:"password" form:"password"`
}

var TIME_TO_EXPIRE = time.Minute * 15

var ID_ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (r *Router) createNewRecovery(id string, userID int) error {
	// verify if one for the user already exists
	recoveries, err := r.recoveryQueries.ReadByUserID(r.ctx, int32(userID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	for _, recovery := range recoveries {
		_, _ = r.recoveryQueries.DeleteByID(r.ctx, recovery.ID)
	}

	createParams := recovery.CreateParams{
		ID:     id,
		UserID: int32(userID),
	}
	_, err = r.recoveryQueries.Create(r.ctx, createParams)
	return err
}

func (r *Router) RequestRecovery(ctx echo.Context) error {
	recoveryDTO := new(RequestRecoveryDTO)
	err := ctx.Bind(recoveryDTO)
	if err != nil {
		errorJSON := NewErrorJson(InvalidForm)
		return ctx.JSON(http.StatusBadRequest, errorJSON)
	}

	if !services.IsValidEmail(recoveryDTO.Email) {
		errorJSON := NewErrorJson("invalid email format")
		return ctx.JSON(http.StatusBadRequest, errorJSON)
	}

	user, err := r.userQueries.ReadUserByEmail(r.ctx, recoveryDTO.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorJSON := NewErrorJson("invalid email")
			return ctx.JSON(http.StatusNotFound, errorJSON)
		}

		errorJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errorJSON)
	}

	recoveryID, err := gonanoid.Generate(ID_ALPHABET, 64)
	if err != nil {
		ctx.Logger().Errorf("could not generate nanoid: %v\n", err)
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	err = r.createNewRecovery(recoveryID, int(user.ID))

	if err != nil {
		ctx.Logger().Errorf("could not generate recovery: %v", err)
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	return ctx.NoContent(http.StatusNoContent)
}
