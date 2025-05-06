package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"wilin/src/database"
)

type ProposalDTO struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Entry    string `json:"entry"`
	Pos      string `json:"pos"`
	Gloss    string `json:"gloss"`
	Notes    string `json:"notes"`
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
