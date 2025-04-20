package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type WordModel struct {
	Id    int
	Entry string
	Pos   string
	Gloss string
	Notes string
}

type Fields struct {
	Entry bool
	Pos   bool
	Gloss bool
	Notes bool
}

func (f *Fields) Map() map[string]bool {
	return map[string]bool{
		"entry": f.Entry,
		"pos":   f.Pos,
		"gloss": f.Gloss,
		"notes": f.Notes,
	}
}

type Column = string

const (
	Entry Column = "entry"
	Pos   Column = "pos"
	Gloss Column = "gloss"
	Notes Column = "notes"
)

type SearchParameters struct {
	Search string
	Fields Fields
	Column Column
	Page   int
}

const PageSize = 100

var (
	ErrNoChange = errors.New("Database was not changed")
)

func GetWordFromRecord(record []string) (WordModel, error) {
	recordLength, word := len(record), WordModel{}
	if recordLength != 3 && recordLength != 4 {
		return word, errors.New("GetWordFromRecord, incorrect record size")
	}

	word.Entry, word.Pos, word.Gloss = record[0], record[1], record[2]
	if recordLength == 4 {
		word.Notes = record[3]
	}
	return word, nil
}

type WordDao struct {
	Db *sql.DB
}

func (dao *WordDao) scanRow(row *sql.Row) (WordModel, error) {
	var word WordModel
	err := row.Scan(&word.Id, &word.Entry, &word.Pos, &word.Gloss, &word.Notes)
	if err != nil {
		return WordModel{}, fmt.Errorf("scanRow, failed at scanning row: %w", err)
	}
	return word, nil
}

func (dao *WordDao) scanRows(rows *sql.Rows) (WordModel, error) {
	var word WordModel
	err := rows.Scan(&word.Id, &word.Entry, &word.Pos, &word.Gloss, &word.Notes)
	if err != nil {
		return WordModel{}, fmt.Errorf("scanRows, failed at scanning row: %w", err)
	}
	return word, nil
}

func (dao *WordDao) CreateWord(word *WordModel) (int64, error) {
	result, err := dao.Db.Exec(`
	INSERT INTO words (entry, pos, gloss, notes)
	VALUES (?, ?, ?, ?)
	`, word.Entry, word.Pos, word.Gloss, word.Notes)
	if err != nil {
		return -1, fmt.Errorf("CreateWord, failed inserting word into database: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("CreateWord, failed at fetching last insert id: %w", err)
	}
	return id, nil
}

func (dao *WordDao) ReadAllWords() ([]WordModel, error) {
	rows, err := dao.Db.Query("SELECT id, entry, pos, gloss, notes FROM words ORDER BY entry")
	if err != nil {
		return nil, fmt.Errorf("ReadAllWords, error querying rows: %w", err)
	}

	var words []WordModel
	for rows.Next() {
		word, err := dao.scanRows(rows)
		if err != nil {
			return words, fmt.Errorf("ReadAllWords, failed at scanning rows: %w", err)
		}
		words = append(words, word)
	}
	return words, nil
}

func (dao *WordDao) generateQueryRestriction(searchParameters SearchParameters) (string, []interface{}) {
	var query string
	var args []interface{}
	searchPattern := "%" + searchParameters.Search + "%"

	for column, enabled := range searchParameters.Fields.Map() {
		if enabled {
			query += fmt.Sprintf("OR %s LIKE ? ", column)
			args = append(args, searchPattern)
		}
	}

	query += fmt.Sprintf("ORDER BY %s ", searchParameters.Column)

	page := searchParameters.Page - 1
	if page < 0 {
		return query, args
	}
	query += fmt.Sprintf("LIMIT %d OFFSET %d", PageSize, PageSize*page)

	return query, args
}

func (dao *WordDao) ReadWordBySearch(parameters SearchParameters) ([]WordModel, int, error) {
	query := "SELECT id, entry, pos, gloss, notes FROM words WHERE 1=0 "
	restriction, args := dao.generateQueryRestriction(parameters)
	query += restriction

	rows, err := dao.Db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ReadWordBySearch, error querying rows: %w", err)
	}

	var words []WordModel
	for rows.Next() {
		word, err := dao.scanRows(rows)
		if err != nil {
			return words, 0, fmt.Errorf("ReadWordBySearch, failed at scanning rows: %w", err)
		}
		words = append(words, word)
	}

	pageCount, err := dao.getPageCount(parameters)
	if err != nil {
		return words, 0, fmt.Errorf("ReadWordBySearch, failed at scanning rows: %w", err)
	}

	return words, pageCount, nil
}

func (dao *WordDao) getPageCount(parameters SearchParameters) (int, error) {
	query := "SELECT COUNT(*) FROM words WHERE 1=0 "
	parameters.Page = 0
	restriction, args := dao.generateQueryRestriction(parameters)
	query += restriction

	row := dao.Db.QueryRow(query, args...)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("getPageCount, failed at scanning row: %w", err)
	}
	return ((count - 1) / PageSize) + 1, nil
}

func (dao *WordDao) ReadWordById(id int) (WordModel, error) {
	row := dao.Db.QueryRow("SELECT id, entry, pos, gloss, notes FROM words WHERE id=?", id)
	word, err := dao.scanRow(row)
	if err != nil {
		return WordModel{}, fmt.Errorf("ReadWordById, failed at scanning row: %w", err)
	}
	return word, nil
}

func (dao *WordDao) DeleteWord(word *WordModel) error {
	return dao.DeleteWordById(word.Id)
}

func (dao *WordDao) DeleteWordById(id int) error {
	result, err := dao.Db.Exec("DELETE FROM words WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("DeleteWordById, failed at deleting word: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteWordById, failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("DeleteWordById, no rows were deleted (row with id %d does not exist): %w", id, sql.ErrNoRows)
	}

	return nil
}

func (dao *WordDao) UpdateWord(word *WordModel) error {
	err := dao.Db.QueryRow("SELECT * FROM words WHERE id=?", word.Id).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("UpdateWord, no rows were updated (row with id %d does not exist): %w", word.Id, sql.ErrNoRows)
	}

	result, err := dao.Db.Exec(`
	UPDATE words
	SET entry=?, pos=?, gloss=?, notes=?
	WHERE id=?
	`, word.Entry, word.Pos, word.Gloss, word.Notes, word.Id)

	if err != nil {
		return fmt.Errorf("UpdateWord, failed at updating word: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateWord, failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("UpdateWord, no rows were updated (row with id %d does not exist): %w", word.Id, ErrNoChange)
	}

	return nil
}
