package extracter

import (
	"encoding/csv"
	"fmt"
	"os"
	"wilin/src/database"
)

func GetRecordsFromCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("GetRecordsFromCsv, failed opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("GetRecordsFromCsv, failed reading csv data: %w", err)
	}
	return records, nil
}

func getWordsFromRecords(records [][]string) ([]database.WordModel, error) {
	var words []database.WordModel
	for _, record := range records {
		word, err := database.GetWordFromRecord(record)
		if err != nil {
			return words, fmt.Errorf("getWordsFromRecords, failed to parse record %s: %w", record, err)
		}
		words = append(words, word)
	}
	return words, nil
}

func AddWordsToDatabase(fileName string) error {
	records, err := GetRecordsFromCsv(fileName)
	if err != nil {
		return fmt.Errorf("main, failed getting csv records: %w", err)
	}

	words, err := getWordsFromRecords(records)
	if err != nil {
		return fmt.Errorf("main, failed at getting words from record: %w", err)
	}

	dao := database.WordDao{}
	for _, word := range words {
		err := dao.CreateWord(&word)
		if err != nil {
			return fmt.Errorf("main, failed to add word \"%s\": %w", word.Entry, err)
		}
	}
	return nil
}
