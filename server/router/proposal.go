package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProposalDTO struct {
	Id       int    `json:"id" form:"id"`
	UserId   int    `json:"userId" form:"userId"`
	Username string `json:"username,omitempty" form:"username"`
	Entry    string `json:"entry" form:"entry"`
	Pos      string `json:"pos" form:"pos"`
	Gloss    string `json:"gloss" form:"gloss"`
	Notes    string `json:"notes" form:"notes"`
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
