package router

import (
	"errors"
	"net/http"
	"slices"

	"wilin.info/api/database/kalan"

	"github.com/labstack/echo/v4"
)

type KalanDTO struct {
	ID    int    `json:"id" form:"id"`
	Entry string `json:"entry" form:"entry"`
	Pos   string `json:"pos" form:"pos"`
	Gloss string `json:"gloss" form:"gloss"`
	Notes string `json:"notes" form:"notes"`
}

func NewKalanDTO(id int, entry string, pos string, gloss string, notes string) KalanDTO {
	return KalanDTO{
		ID:    id,
		Entry: entry,
		Pos:   pos,
		Gloss: gloss,
		Notes: notes,
	}
}

func validateKalanJson(kalan *KalanDTO) error {
	if kalan.Entry == "" {
		return errNoEntry
	}
	if kalan.Pos == "" {
		return errNoPos
	}
	if kalan.Gloss == "" {
		return errNoGloss
	}
	if kalan.ID == 0 {
		return errNoId
	}
	return nil
}

type KalanArrayDTO struct {
	Kalans     []KalanDTO `json:"kalans"`
	Page       int        `json:"page"`
	KalanCount int        `json:"kalanCount"`
	PageCount  int        `json:"pageCount"`
}

func (arr *KalanArrayDTO) AddKalan(kalan KalanDTO) {
	arr.Kalans = append(arr.Kalans, kalan)
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

// Define the Handlers for the kalan related routes

func (r *Router) GetAllKalan(ctx echo.Context) error {
	var kalanArrayDTO KalanArrayDTO

	kalans, err := r.kalanQueries.ReadKalan(r.ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorJson("Failed to fetch words"))
	}

	for _, kalan := range kalans {
		kalanDTO := NewKalanDTO(int(kalan.ID), kalan.Entry, kalan.Pos, kalan.Gloss, kalan.Notes)
		kalanArrayDTO.AddKalan(kalanDTO)
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
		kalanArrayDTO.AddKalan(kalanDTO)
	}

	searchCountParams := kalan.ReadKalanSearchCountParams{
		Search:  searchParams.Search,
		Isentry: searchParams.Isentry,
		Ispos:   searchParams.Ispos,
		Isgloss: searchParams.Isgloss,
		Isnotes: searchParams.Isnotes,
	}
	kalanCount, err := r.kalanQueries.ReadKalanSearchCount(r.ctx, searchCountParams)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorJson("Could not fetch words"))
	}

	pageCount := getPageCount(int(kalanCount), PAGE_SIZE)

	kalanArrayDTO.Page = searchQueryDTO.Page
	kalanArrayDTO.KalanCount = int(kalanCount)
	kalanArrayDTO.PageCount = pageCount

	return ctx.JSON(http.StatusOK, kalanArrayDTO)
}

func (r *Router) AddKalan(ctx echo.Context) error {
	var kalanDTO KalanDTO
	err := ctx.Bind(&kalanDTO)
	if err != nil {
		errJSON := NewErrorJson(errInvalidFormat.Error())
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	err = validateKalanJson(&kalanDTO)
	if err != nil && !errors.Is(err, errNoId) {
		errJSON := NewErrorJson(err.Error())
		return ctx.JSON(http.StatusBadRequest, errJSON)
	}

	createParams := kalan.CreateKalanParams{
		Entry: kalanDTO.Entry,
		Pos:   kalanDTO.Pos,
		Gloss: kalanDTO.Gloss,
		Notes: kalanDTO.Notes,
	}

	result, err := r.kalanQueries.CreateKalan(r.ctx, createParams)
	if err != nil {
		errJSON := NewErrorJson("could not add kalan to database")
		return ctx.JSON(http.StatusInternalServerError, errJSON)
	}

	kalanID, err := result.LastInsertId()
	if err != nil {
		ctx.Logger().Errorf("error fetching result id: %v\n", err.Error())
	}

	kalanDTO.ID = int(kalanID)
	return ctx.JSON(http.StatusCreated, kalanDTO)
}
