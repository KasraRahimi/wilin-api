package router

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"wilin.info/api/database/recovery"
	"wilin.info/api/database/users"
	"wilin.info/api/server/services"
)

type RequestRecoveryDTO struct {
	Email string `json:"email" form:"email"`
}

type ChangePasswordDTO struct {
	Password string `json:"password" form:"password"`
	ID       string `param:"id"`
}

func validatePassword(p string) error {
	if len(p) < MIN_PASSWORD_LENGTH {
		return errors.New(PasswordTooShort)
	}
	if len(p) > MAX_PASSWORD_LENGTH {
		return errors.New(PasswordTooLong)
	}
	return nil
}

var TIME_TO_EXPIRE = time.Minute * 15

var ID_ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func isExpired(t *time.Time) bool {
	expiration := t.Add(TIME_TO_EXPIRE)
	return expiration.Before(time.Now())
}

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

func (r *Router) ChangePassword(ctx echo.Context) error {
	newPasswordDTO := new(ChangePasswordDTO)
	err := ctx.Bind(newPasswordDTO)
	if err != nil {
		errJSON := NewErrorJson(InvalidForm)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	recovery, err := r.recoveryQueries.ReadByID(r.ctx, newPasswordDTO.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errJSON := NewErrorJson("invalid recovery id")
			return ctx.JSON(http.StatusNotFound, errJSON)
		}

		ctx.Logger().Errorf("could not find recovery item: %v\n", err)
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	if isExpired(&recovery.CreatedAt) {
		ctx.Logger().Infof("recovery id is expired. deleting %v...\n", recovery.ID)
		_, _ = r.recoveryQueries.DeleteByID(r.ctx, recovery.ID)
		errJSON := NewErrorJson("invalid recovery id")
		return ctx.JSON(http.StatusNotFound, errJSON)
	}

	err = validatePassword(newPasswordDTO.Password)
	if err != nil {
		errJSON := NewErrorJson(err.Error())
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	passwordHash, err := services.GeneratePasswordHash(newPasswordDTO.Password)
	if err != nil {
		ctx.Logger().Errorf("could not generate password hash: %v\n", err)
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	updateParams := users.UpdatePasswordParams{
		ID:       recovery.UserID,
		Password: passwordHash,
	}
	_, err = r.userQueries.UpdatePassword(r.ctx, updateParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errJSON := NewErrorJson("invalid recovery id")
			return ctx.JSON(http.StatusNotFound, errJSON)
		}

		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	_, _ = r.recoveryQueries.DeleteByID(r.ctx, recovery.ID)

	return ctx.NoContent(http.StatusNoContent)
}
