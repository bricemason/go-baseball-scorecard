package main

import (
	"encoding/csv"
	"strings"
)

// GetRecords Given a slice of strings representing comma-delimited data, returns a multidimensional slice of records
func GetRecords(source []string) [][]string {
	csvSource := strings.Join(source, "\r\n")
	csvReader := csv.NewReader(strings.NewReader(csvSource))

	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}
