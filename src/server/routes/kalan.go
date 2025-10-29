package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"wilin.com/api/src/database"

	"github.com/gin-gonic/gin"
)

type WilinWordJson struct {
	ID    int    `json:"id,omitempty" form:"id"`
	Entry string `json:"entry" form:"entry"`
	Pos   string `json:"pos" form:"pos"`
	Gloss string `json:"gloss" form:"gloss"`
	Notes string `json:"notes" form:"notes"`
}

func getJsonFromWordModel(word *database.WordModel) WilinWordJson {
	return WilinWordJson{word.Id, word.Entry, word.Pos, word.Gloss, word.Notes}
}

func getWordModelFromJson(word *WilinWordJson) database.WordModel {
	return database.WordModel{
		Id:    word.ID,
		Entry: word.Entry,
		Pos:   word.Pos,
		Gloss: word.Gloss,
		Notes: word.Notes,
	}
}

func (s *Server) getSearchParameters(ctx *gin.Context) *database.SearchParameters {
	parameters := database.SearchParameters{}
	search, isSearch := ctx.GetQuery("search")
	fieldString, isFieldStrings := ctx.GetQuery("fields")
	column, isColumn := ctx.GetQuery("sort")
	pageString, isPageString := ctx.GetQuery("page")

	if !(isSearch || isColumn || isFieldStrings || isPageString) {
		return nil
	}

	parameters.Search = search

	if !isFieldStrings {
		fieldString = "entry,pos,gloss,notes"
	}

	fieldStrings := strings.Split(fieldString, ",")
	fields := database.Fields{}
	for _, field := range fieldStrings {
		switch strings.ToLower(field) {
		case "entry":
			fields.Entry = true
		case "pos":
			fields.Pos = true
		case "gloss":
			fields.Gloss = true
		case "notes":
			fields.Notes = true
		}
	}

	parameters.Fields = fields

	if !isColumn {
		column = "entry"
	}

	switch column {
	case "entry":
		parameters.Column = database.Entry
	case "pos":
		parameters.Column = database.Pos
	case "gloss":
		parameters.Column = database.Gloss
	case "notes":
		parameters.Column = database.Notes
	default:
		parameters.Column = database.Entry
	}

	parameters.Page, _ = strconv.Atoi(pageString)
	if parameters.Page < 1 {
		parameters.Page = 1
	}

	return &parameters
}

func (s *Server) HandleGetKalan(ctx *gin.Context) {
	words, err := s.WordDao.ReadAllWords()

	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	var wordsJson []WilinWordJson
	for _, word := range words {
		wordsJson = append(wordsJson, getJsonFromWordModel(&word))
	}

	ctx.JSON(http.StatusOK, wordsJson)
}

func (s *Server) HandleGetKalanPaginated(ctx *gin.Context) {
	parameters := s.getSearchParameters(ctx)
	words, pageCount, err := s.WordDao.ReadWordBySearch(*parameters)

	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	var wordsJson []WilinWordJson
	for _, word := range words {
		wordsJson = append(wordsJson, getJsonFromWordModel(&word))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pageCount": pageCount,
		"words":     wordsJson,
	})
}

func (s *Server) HandleGetKalanById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Badly formatted id")
		return
	}

	word, err := s.WordDao.ReadWordById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Error(err)
			ctx.String(http.StatusNotFound, "Word with id=%d not found", id)
			return
		}
		ctx.Error(err)
		ctx.String(http.StatusInternalServerError, "Something went wrong fetching the word")
		return
	}
	ctx.JSON(http.StatusOK, getJsonFromWordModel(&word))
}

func (s *Server) validateKalanJson(kalan *WilinWordJson) error {
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

func (s *Server) HandlePostKalan(ctx *gin.Context) {
	var kalanJson WilinWordJson
	if err := ctx.ShouldBind(&kalanJson); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(errInvalidFormat.Error()))
		return
	}

	if err := s.validateKalanJson(&kalanJson); err != nil && !errors.Is(err, errNoId) {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(err.Error()))
		return
	}

	word := getWordModelFromJson(&kalanJson)
	id, err := s.WordDao.CreateWord(&word)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	kalanJson.ID = int(id)
	ctx.JSON(http.StatusCreated, kalanJson)
}

func (s *Server) HandleDeleteKalan(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Badly formatted id")
		return
	}

	err = s.WordDao.DeleteWordById(id)
	if err != nil {
		ctx.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.String(http.StatusNotFound, "Cannot find word with id=%d", id)
			return
		}
		ctx.String(http.StatusInternalServerError, "Something went wrong deleting the word")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) HandlePutKalan(ctx *gin.Context) {
	var kalanJson WilinWordJson
	if err := ctx.ShouldBind(&kalanJson); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(errInvalidFormat.Error()))
		return
	}

	if err := s.validateKalanJson(&kalanJson); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, GetErrorJson(err.Error()))
		return
	}

	word := getWordModelFromJson(&kalanJson)
	err := s.WordDao.UpdateWord(&word)
	if err != nil {
		ctx.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.String(http.StatusNotFound, "Cannot find word id=%d, entry=%s", word.Id, word.Entry)
			return
		}
		if errors.Is(err, database.ErrNoChange) {
			ctx.String(http.StatusOK, "Update did not change word")
			return
		}

		ctx.JSON(http.StatusInternalServerError, GetErrorJson("something went wrong"))
		return
	}

	ctx.Status(http.StatusCreated)
}
