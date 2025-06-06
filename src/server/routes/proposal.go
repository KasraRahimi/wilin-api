package routes

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wilin/src/database"
	"wilin/src/database/permissions"
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

	if len(proposalsDTO) == 0 {
		ctx.JSON(http.StatusNotFound, []ProposalDTO{})
		return
	}
	ctx.JSON(http.StatusOK, proposalsDTO)
}

func (s *Server) HandleGetProposalById(ctx *gin.Context) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson(errNoUserFromCtx.Error()))
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, GetErrorJson("unauthorized"))
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, GetErrorJson("invalid id"))
		return
	}

	proposal, err := s.ProposalDao.ReadProposalById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	if !permissions.CanRolePermission(user.Role, permissions.VIEW_ALL_PROPOSAL) {
		canUserViewSelfProposal := permissions.CanRolePermission(user.Role, permissions.VIEW_SELF_PROPOSAL)
		isOwner := user.Id == proposal.UserId
		if !canUserViewSelfProposal || !isOwner {
			ctx.JSON(http.StatusForbidden, GetErrorJson("permission denied"))
			return
		}
	}

	proposalDTO := NewProposalDTOFromModel(proposal, "")
	ctx.JSON(http.StatusOK, proposalDTO)
}

func (s *Server) HandleDeleteProposal(ctx *gin.Context) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson(errNoUserFromCtx.Error()))
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, GetErrorJson("unauthorized"))
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, GetErrorJson("invalid id"))
		return
	}

	proposal, err := s.ProposalDao.ReadProposalById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	if !permissions.CanRolePermission(user.Role, permissions.DELETE_ALL_PROPOSAL) {
		canUserDeleteSelfProposal := permissions.CanRolePermission(user.Role, permissions.DELETE_SELF_PROPOSAL)
		isOwner := user.Id == proposal.UserId
		if !canUserDeleteSelfProposal || !isOwner {
			ctx.JSON(http.StatusForbidden, GetErrorJson("permission denied"))
			return
		}
	}

	err = s.ProposalDao.Delete(proposal)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) HandlePutProposal(ctx *gin.Context) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson(errNoUserFromCtx.Error()))
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, GetErrorJson("unauthorized"))
		return
	}

	var proposalDTO ProposalDTO
	if err := ctx.ShouldBind(&proposalDTO); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(errInvalidFormat.Error()))
		return
	}
	if err = s.validateProposalJSON(&proposalDTO); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(err.Error()))
		return
	}

	originalProposal, err := s.ProposalDao.ReadProposalById(proposalDTO.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	if !permissions.CanRolePermission(user.Role, permissions.MODIFY_ALL_PROPOSAL) {
		canUserDeleteSelfProposal := permissions.CanRolePermission(user.Role, permissions.DELETE_SELF_PROPOSAL)
		isOwner := user.Id == originalProposal.UserId
		if !canUserDeleteSelfProposal || !isOwner {
			ctx.JSON(http.StatusForbidden, GetErrorJson("permission denied"))
			return
		}
	}

	proposalDTO.UserId = originalProposal.UserId
	newProposal := proposalDTO.ToModel()
	err = s.ProposalDao.Update(&newProposal)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	newProposalDTO := NewProposalDTOFromModel(&newProposal, "")
	ctx.JSON(http.StatusOK, newProposalDTO)
}
