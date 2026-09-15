package converter

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Employee struct {
	ID            string
	LastName      string
	FirstName     string
	Email         string
	LastNameKana  string
	FirstNameKana string
	Note          string
}

type DictionaryEntry struct {
	Reading      string
	Word         string
	PartOfSpeech string
	Comment      string
}

func ParseGoogleTSV(r io.Reader) ([]DictionaryEntry, error) {
	reader := csv.NewReader(r)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	var entries []DictionaryEntry

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 3 || len(record) > 4 {
			return nil, fmt.Errorf("expected 3 or 4 columns, got %d", len(record))
		}

		entry := DictionaryEntry{
			Reading:      record[0],
			Word:         record[1],
			PartOfSpeech: record[2],
		}
		if len(record) == 4 {
			entry.Comment = record[3]
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func ParseCSV(r io.Reader) ([]Employee, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	var employees []Employee

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 6 || len(record) > 7 {
			return nil, fmt.Errorf("expected 6 or 7 columns, got %d", len(record))
		}
		e := Employee{
			ID:            record[0],
			LastName:      record[1],
			FirstName:     record[2],
			Email:         record[3],
			LastNameKana:  record[4],
			FirstNameKana: record[5],
		}
		if len(record) == 7 {
			e.Note = record[6]
		}
		employees = append(employees, e)
	}

	return employees, nil
}
