package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"wilin/src/database"

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

func (s *Server) getSearchAndFields(ctx *gin.Context) (string, database.Fields) {
	search, _ := ctx.GetQuery("search")
	fieldString, isFieldStrings := ctx.GetQuery("fields")
	if !isFieldStrings {
		return search, database.Fields{Entry: true, Pos: true, Gloss: true, Notes: true}
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
	return search, fields
}

func (s *Server) HandleGetKalan(ctx *gin.Context) {
	var words []database.WordModel
	var err error
	search, fields := s.getSearchAndFields(ctx)
	if search == "" {
		words, err = s.WordDao.ReadAllWords()

		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else {
		words, err = s.WordDao.ReadWordBySearch(search, fields)

		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	var wordsJson []WilinWordJson
	for _, word := range words {
		wordsJson = append(wordsJson, getJsonFromWordModel(&word))
	}
	ctx.JSON(http.StatusOK, wordsJson)
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

func (s *Server) HandlePostKalan(ctx *gin.Context) {
	var kalanJson WilinWordJson
	if err := ctx.ShouldBind(&kalanJson); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusBadRequest, "Incorrectly formatted")
		return
	}

	word := getWordModelFromJson(&kalanJson)
	id, err := s.WordDao.CreateWord(&word)
	if err != nil {
		ctx.Error(err)
		ctx.String(http.StatusInternalServerError, "Something went wrong adding this word")
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

	ctx.Status(http.StatusOK)
}

func (s *Server) HandlePutKalan(ctx *gin.Context) {
	var kalanJson WilinWordJson
	if err := ctx.BindJSON(&kalanJson); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusBadRequest, "Incorrectly formatted")
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

		ctx.String(http.StatusInternalServerError, "Something went wrong updating the word")
		return
	}

	ctx.Status(http.StatusCreated)
}
