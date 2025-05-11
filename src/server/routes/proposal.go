package routes

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"wilin/src/database"
)

type ProposalDTO struct {
	Id       int    `json:"id" form:"id"`
	UserId   int    `json:"userId" form:"userId"`
	Username string `json:"username" form:"username"`
	Entry    string `json:"entry" form:"entry"`
	Pos      string `json:"pos" form:"pos"`
	Gloss    string `json:"gloss" form:"gloss"`
	Notes    string `json:"notes" form:"notes"`
}

func (p *ProposalDTO) ToModel() database.ProposalModel {
	return database.ProposalModel{
		Id:     p.Id,
		UserId: p.UserId,
		Entry:  p.Entry,
		Pos:    p.Pos,
		Gloss:  p.Gloss,
		Notes:  p.Notes,
	}
}

func NewProposalDTOFromUsernameModel(model *database.ProposalUsernameModel) ProposalDTO {
	return ProposalDTO{
		Id:       model.Id,
		UserId:   model.UserId,
		Username: model.Username,
		Entry:    model.Entry,
		Pos:      model.Pos,
		Gloss:    model.Gloss,
		Notes:    model.Notes,
	}
}

func NewProposalDTOFromModel(model *database.ProposalModel, username string) ProposalDTO {
	return ProposalDTO{
		Id:       model.Id,
		UserId:   model.UserId,
		Username: username,
		Entry:    model.Entry,
		Pos:      model.Pos,
		Gloss:    model.Gloss,
		Notes:    model.Notes,
	}
}

func (s *Server) validateProposalJSON(dto *ProposalDTO) error {
	if dto.Entry == "" {
		return errNoEntry
	}
	if dto.Pos == "" {
		return errNoPos
	}
	if dto.Gloss == "" {
		return errNoGloss
	}
	if dto.UserId == 0 {
		return errNoUserID
	}
	if dto.Id == 0 {
		return errNoId
	}
	return nil
}

func (s *Server) HandlePostProposal(ctx *gin.Context) {
	var proposal ProposalDTO
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, GetErrorJson("something went wrong"))
		return
	}

	if err := ctx.ShouldBind(&proposal); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(errInvalidFormat.Error()))
		return
	}

	proposal.UserId = user.Id

	if err := s.validateProposalJSON(&proposal); err != nil && !errors.Is(err, errNoId) {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(err.Error()))
		return
	}

	proposalModel := proposal.ToModel()
	id, err := s.ProposalDao.CreateProposal(&proposalModel)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	proposal.Id = int(id)
	ctx.JSON(http.StatusCreated, proposal)
}

func (s *Server) HandleGetAllProposals(ctx *gin.Context) {
	proposalModels, err := s.ProposalDao.ReadAllProposalsWithUsername()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong fetching all proposals"))
		return
	}
	var proposalsDTO []ProposalDTO
	for _, model := range proposalModels {
		proposalsDTO = append(proposalsDTO, NewProposalDTOFromUsernameModel(&model))
	}
	ctx.JSON(http.StatusOK, proposalsDTO)
}

func (s *Server) HandleGetMyProposals(ctx *gin.Context) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("could not fetch user"))
		return
	}
	proposals, err := s.ProposalDao.ReadProposalsByUserId(user.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson(err.Error()))
		return
	}
	var proposalsDTO []ProposalDTO
	for _, model := range proposals {
		proposal := NewProposalDTOFromModel(&model, user.Username)
		proposalsDTO = append(proposalsDTO, proposal)
	}
	ctx.JSON(http.StatusOK, proposalsDTO)
}
