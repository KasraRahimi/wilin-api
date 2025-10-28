package router

import (
	"net/http"
	"slices"
	"wilin/database/kalan"

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
	Kalans     []KalanDTO `json:"kalans"`
	Page       int        `json:"page"`
	KalanCount int        `json:"kalanCount"`
	PageCount  int        `json:"pageCount"`
}

type SearchQueryDTO struct {
	Search string `query:"search"`
	Fields string `query:"fields"`
	Sort   string `query:"sort"`
	Page   int    `query:"page"`
}

type Fields struct {
	IsEntry bool
	IsPos   bool
	IsGloss bool
	IsNotes bool
}

const PAGE_SIZE = 100

func getPageCount(kalanCount int, pageSize int) int {
	return ((kalanCount - 1) / pageSize) + 1
}

func NewFields(fieldsArray []string) Fields {
	if len(fieldsArray) < 1 {
		return Fields{
			IsEntry: true,
			IsPos:   true,
			IsGloss: true,
			IsNotes: true,
		}
	}
	var fields Fields
	fields.IsEntry = slices.Contains(fieldsArray, "entry")
	fields.IsPos = slices.Contains(fieldsArray, "pos")
	fields.IsGloss = slices.Contains(fieldsArray, "gloss")
	fields.IsNotes = slices.Contains(fieldsArray, "notes")
	return fields
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

func (r *Router) GetKalanBySearch(ctx echo.Context) error {
	var searchQueryDTO SearchQueryDTO
	err := ctx.Bind(&searchQueryDTO)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorJson("Invalid search"))
	}

	fields := NewFields(splitQuery(searchQueryDTO.Fields))

	searchParams := kalan.ReadKalanBySearchParams{
		Search:  searchQueryDTO.Search,
		Isentry: fields.IsEntry,
		Ispos:   fields.IsPos,
		Isgloss: fields.IsGloss,
		Isnotes: fields.IsNotes,
		Sort:    searchQueryDTO.Sort,
		Limit:   PAGE_SIZE,
		Offset:  int32(PAGE_SIZE * searchQueryDTO.Page),
	}

	kalans, err := r.kalanQueries.ReadKalanBySearch(r.ctx, searchParams)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorJson("Could not fetch words"))
	}

	var kalanArrayDTO KalanArrayDTO
	for _, kalan := range kalans {
		kalanDTO := NewKalanDTO(int(kalan.ID), kalan.Entry, kalan.Pos, kalan.Gloss, kalan.Notes)
		kalanArrayDTO.Kalans = append(kalanArrayDTO.Kalans, kalanDTO)
	}

	kalanCount, err := r.kalanQueries.ReadKalanCount(r.ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorJson("Could not fetch words"))
	}

	pageCount := getPageCount(int(kalanCount), PAGE_SIZE)

	kalanArrayDTO.Page = searchQueryDTO.Page
	kalanArrayDTO.KalanCount = int(kalanCount)
	kalanArrayDTO.PageCount = pageCount

	return ctx.JSON(http.StatusOK, kalanArrayDTO)
}
