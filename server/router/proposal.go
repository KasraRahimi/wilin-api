package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"wilin.info/api/database/proposal"
	"wilin.info/api/server/services"
)

type ProposalIDDTO struct {
	ID int `param:"id"`
}

type ProposalDTO struct {
	Id       int    `json:"id" form:"id"`
	UserId   int    `json:"userId" form:"userId"`
	Username string `json:"username,omitempty" form:"username"`
	Entry    string `json:"entry" form:"entry"`
	Pos      string `json:"pos" form:"pos"`
	Gloss    string `json:"gloss" form:"gloss"`
	Notes    string `json:"notes" form:"notes"`
}

func validateProposalJSON(dto *ProposalDTO) error {
	if dto.Entry == "" {
		return ErrNoEntry
	}
	if dto.Pos == "" {
		return ErrNoPos
	}
	if dto.Gloss == "" {
		return ErrNoGloss
	}
	if dto.Id == 0 {
		return ErrNoId
	}
	return nil
}

type ProposalArrDTO struct {
	Proposals []ProposalDTO `json:"proposals"`
}

func (pArr *ProposalArrDTO) AddProposal(p ProposalDTO) {
	pArr.Proposals = append(pArr.Proposals, p)
}

func (r *Router) GetAllProposals(ctx echo.Context) error {
	proposals, err := r.proposalQueries.ReadAllProposalsWithUsername(r.ctx)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposalArrDTO := new(ProposalArrDTO)

	for _, p := range proposals {
		proposalDto := ProposalDTO{
			Id:       int(p.ID),
			UserId:   int(p.UserID.Int32),
			Username: p.Username,
			Entry:    p.Entry,
			Pos:      p.Pos,
			Gloss:    p.Gloss,
			Notes:    p.Notes,
		}
		proposalArrDTO.AddProposal(proposalDto)
	}

	return ctx.JSON(http.StatusOK, *proposalArrDTO)
}

func (r *Router) PostProposal(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(int)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	proposalDTO := new(ProposalDTO)
	err := ctx.Bind(proposalDTO)
	if err != nil {
		errJSON := NewErrorJson(InvalidForm)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	err = validateProposalJSON(proposalDTO)
	if err != nil && !errors.Is(err, ErrNoId) {
		errJSON := NewErrorJson(err.Error())
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	createParams := proposal.CreateProposalParams{
		UserID: sql.NullInt32{Int32: int32(userID), Valid: true},
		Entry:  proposalDTO.Entry,
		Pos:    proposalDTO.Pos,
		Gloss:  proposalDTO.Gloss,
		Notes:  proposalDTO.Notes,
	}
	result, err := r.proposalQueries.CreateProposal(r.ctx, createParams)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposalID, err := result.LastInsertId()
	if err != nil {
		ctx.Logger().Errorf("error fetching result id: %v\n", err.Error())
	}

	proposalDTO.Id = int(proposalID)
	proposalDTO.UserId = userID

	return ctx.JSON(http.StatusCreated, proposalDTO)
}

func (r *Router) GetMyProposals(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(int)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	proposals, err := r.proposalQueries.ReadProposalsByUserIDWithUsername(r.ctx, int32(userID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposalArrDTO := ProposalArrDTO{Proposals: []ProposalDTO{}}
	for _, proposal := range proposals {
		proposalDTO := ProposalDTO{
			Id:       int(proposal.ID),
			UserId:   int(proposal.UserID.Int32),
			Username: proposal.Username,
			Entry:    proposal.Entry,
			Pos:      proposal.Pos,
			Gloss:    proposal.Gloss,
			Notes:    proposal.Notes,
		}
		proposalArrDTO.AddProposal(proposalDTO)
	}

	return ctx.JSON(http.StatusOK, proposalArrDTO)
}

func (r *Router) GetProposalByID(ctx echo.Context) error {
	params := ProposalIDDTO{}
	err := ctx.Bind(&params)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	userID, ok := ctx.Get("userID").(int)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	user, err := r.userQueries.ReadUserByID(r.ctx, int32(userID))
	if err != nil {
		ctx.Logger().Errorf("could not fetch user: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposal, err := r.proposalQueries.ReadProposalByIDWithUsername(r.ctx, int32(params.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("no proposal with id=%v", params.ID)
			errJSON := NewErrorJson(msg)
			return ctx.JSON(http.StatusNotFound, errJSON)
		}
		ctx.Logger().Errorf("could not fetch proposal: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userRole := services.NewRole(user.Role)
	isUserOwner := userRole.Can(services.PERMISSION_VIEW_SELF_PROPOSAL) && user.ID == proposal.UserID.Int32
	isUserAdmin := userRole.Can(services.PERMISSION_VIEW_ALL_PROPOSAL)

	if !isUserOwner && !isUserAdmin {
		return ctx.NoContent(http.StatusForbidden)
	}

	proposalDTO := ProposalDTO{
		Id:       int(proposal.ID),
		UserId:   int(proposal.UserID.Int32),
		Username: proposal.Username,
		Entry:    proposal.Entry,
		Pos:      proposal.Pos,
		Gloss:    proposal.Gloss,
		Notes:    proposal.Notes,
	}

	return ctx.JSON(http.StatusOK, proposalDTO)
}

func (r *Router) UpdateProposal(ctx echo.Context) error {
	proposalDTO := ProposalDTO{}
	err := ctx.Bind(&proposalDTO)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	userID, ok := ctx.Get("userID").(int)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	user, err := r.userQueries.ReadUserByID(r.ctx, int32(userID))
	if err != nil {
		ctx.Logger().Errorf("could not fetch user: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	prop, err := r.proposalQueries.ReadProposalByIDWithUsername(r.ctx, int32(proposalDTO.Id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("no proposal with id=%v", proposalDTO.Id)
			errJSON := NewErrorJson(msg)
			return ctx.JSON(http.StatusNotFound, errJSON)
		}
		ctx.Logger().Errorf("could not fetch proposal: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userRole := services.NewRole(user.Role)
	isUserOwner := userRole.Can(services.PERMISSION_MODIFY_SELF_PROPOSAL) && user.ID == prop.UserID.Int32
	isUserAdmin := userRole.Can(services.PERMISSION_MODIFY_ALL_PROPOSAL)

	if !isUserOwner && !isUserAdmin {
		return ctx.NoContent(http.StatusForbidden)
	}

	updateParams := proposal.UpdateParams{
		UserID: sql.NullInt32{Int32: prop.UserID.Int32, Valid: true},
		Entry:  proposalDTO.Entry,
		Pos:    proposalDTO.Pos,
		Gloss:  proposalDTO.Gloss,
		Notes:  proposalDTO.Notes,
		ID:     int32(proposalDTO.Id),
	}

	_, err = r.proposalQueries.Update(r.ctx, updateParams)
	if err != nil {
		ctx.Logger().Errorf("could not update proposal: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposalDTO.UserId = int(prop.UserID.Int32)
	proposalDTO.Username = prop.Username

	return ctx.JSON(http.StatusOK, proposalDTO)
}

func (r *Router) DeleteProposal(ctx echo.Context) error {
	params := ProposalIDDTO{}
	err := ctx.Bind(&params)
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	userID, ok := ctx.Get("userID").(int)
	if !ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	user, err := r.userQueries.ReadUserByID(r.ctx, int32(userID))
	if err != nil {
		ctx.Logger().Errorf("could not fetch user: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	proposal, err := r.proposalQueries.ReadProposalByIDWithUsername(r.ctx, int32(params.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("no proposal with id=%v", params.ID)
			errJSON := NewErrorJson(msg)
			return ctx.JSON(http.StatusNotFound, errJSON)
		}
		ctx.Logger().Errorf("could not fetch proposal: %v", err.Error())
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	userRole := services.NewRole(user.Role)
	isUserOwner := userRole.Can(services.PERMISSION_DELETE_SELF_PROPOSAL) && user.ID == proposal.UserID.Int32
	isUserAdmin := userRole.Can(services.PERMISSION_DELETE_ALL_PROPOSAL)

	if !isUserOwner && !isUserAdmin {
		return ctx.NoContent(http.StatusForbidden)
	}

	_, err = r.proposalQueries.Delete(r.ctx, int32(params.ID))
	if err != nil {
		errJSON := NewErrorJson(ServerError)
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	return ctx.NoContent(http.StatusNoContent)
}
