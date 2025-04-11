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

type WordDao struct{}

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
	conn, err := GetConnection()
	if err != nil {
		return -1, fmt.Errorf("CreateWord, failed at getting database connection: %w", err)
	}
	defer conn.Close()

	result, err := conn.Exec(`
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
	conn, err := GetConnection()
	if err != nil {
		return nil, fmt.Errorf("ReadAllWords, failed at getting database connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id, entry, pos, gloss, notes FROM words")
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

func (dao *WordDao) ReadWordById(id int) (WordModel, error) {
	conn, err := GetConnection()
	if err != nil {
		return WordModel{}, fmt.Errorf("ReadWordById, failed at getting database connection: %w", err)
	}
	defer conn.Close()

	row := conn.QueryRow("SELECT id, entry, pos, gloss, notes FROM words WHERE id=?", id)
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
	conn, err := GetConnection()
	if err != nil {
		return fmt.Errorf("DeleteWordById, failed at getting database connection: %w", err)
	}

	result, err := conn.Exec("DELETE FROM words WHERE id=?", id)
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
	conn, err := GetConnection()
	if err != nil {
		return fmt.Errorf("UpdateWord, failed at getting database connection: %w", err)
	}

	err = conn.QueryRow("SELECT * FROM words WHERE id=?", word.Id).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("UpdateWord, no rows were updated (row with id %d does not exist): %w", word.Id, sql.ErrNoRows)
	}

	result, err := conn.Exec(`
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
