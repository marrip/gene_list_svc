package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func readTsv(path string) (data [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	data, err = reader.ReadAll()
	if err != nil {
		return
	}
	return
}

func writeTsv(path string) (tsv *csv.Writer, err error) {
	tsvFile, err := os.Create(path)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not create tsv file %s", path))
		return
	}
	defer tsvFile.Close()
	tsv = csv.NewWriter(tsvFile)
	tsv.Comma = '\t'
	defer tsv.Flush()
	return
}
