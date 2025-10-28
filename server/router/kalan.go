package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type KalanDTO struct {
	Id    int    `json:"id"`
	Entry string `json:"entry"`
	Pos   string `json:"pos"`
	Gloss string `json:"gloss"`
	Notes string `json:"notes"`
}

func NewKalanDTO(id int, entry string, pos string, gloss string, notes string) KalanDTO {
	return KalanDTO{
		Id:    id,
		Entry: entry,
		Pos:   pos,
		Gloss: gloss,
		Notes: notes,
	}
}

type KalanArrayDTO struct {
	Kalans []KalanDTO `json:"kalans"`
}

func (r *Router) GetAllKalan(ctx echo.Context) error {
	var kalanArrayDTO KalanArrayDTO

	kalans, err := r.kalanQueries.ReadKalan(r.ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorJson("Failed to fetch words"))
	}

	for _, kalan := range kalans {
		kalanDTO := NewKalanDTO(int(kalan.ID), kalan.Entry, kalan.Pos, kalan.Gloss, kalan.Notes)
		kalanArrayDTO.Kalans = append(kalanArrayDTO.Kalans, kalanDTO)
	}

	return ctx.JSON(http.StatusOK, kalanArrayDTO)
}
