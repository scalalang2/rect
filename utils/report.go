package utils

import (
	"encoding/csv"
	"os"
)

func ReportToCSV(filename string, data [][]string) {
	file, err := os.Create(filename)
	CheckError(err, "[report-ReportToCSV] failed to open csv file")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, el := range data {
		err := writer.Write(el)
		CheckError(err, "[report-ReportToCSV] failed to write csv file")
	}
}