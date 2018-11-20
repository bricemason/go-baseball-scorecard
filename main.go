package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// hard-coding in 1986 Mets home games for now
	eventData, err := ioutil.ReadFile("data/1986NYN.EVN")
	if err != nil {
		fmt.Println("Error reading event file", err)
		return
	}

	// []string of full records in the file
	records := strings.Split(string(eventData), "\n")

	for _, record := range records {
		// []string of fields in the record
		fields := strings.Split(record, ",")

		// the first item is the field type (id, info, play etc)
		recordType := fields[0]

		if recordType == "id" {
			fmt.Println("Game:", fields[1])
		}
	}
}
