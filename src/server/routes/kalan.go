package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"wilin/src/database"

	"github.com/gin-gonic/gin"
)

type WilinWordJson struct {
	ID    int    `json:"id,omitempty"`
	Entry string `json:"entry"`
	Pos   string `json:"pos"`
	Gloss string `json:"gloss"`
	Notes string `json:"notes"`
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
	if err := ctx.BindJSON(&kalanJson); err != nil {
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
