package main

import (
	"encoding/csv"
	"fmt"
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

// ToInterface takes a slice of strings and returns them as a slice of interface{} to make them spreadable as args to sprintf
func ToInterface(s []string) []interface{} {
	list := make([]interface{}, len(s))
	for i, v := range s {
		list[i] = v
	}

	return list
}

// Contains will check if a string value is contained within a slice of strings
func Contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

// TranslateModifiers takes in a slice of modifier codes and returns the descriptions as a comma-delimited list ready for printing
func TranslateModifiers(modifiers []string) string {
	if len(modifiers) > 0 {
		modifierDescriptions := make([]string, 0)

		for _, modifier := range modifiers {
			modifierDescriptions = append(modifierDescriptions, PlayModifiers[modifier])
		}
		return fmt.Sprintf("(%s)", strings.Join(modifierDescriptions, ","))
	}

	return ""
}

// AppendUnique will add an item to a slice of strings if it does not exist
func AppendUnique(list []string, value string) []string {
	if !Contains(list, value) {
		list = append(list, value)
	}

	return list
}
