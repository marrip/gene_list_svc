package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func (e *Entity) tsvToEntity(tsv []string) (err error) {
	if len(tsv) != 4 {
		err = errors.New("Input tsv needs to have 4 columns")
	}
	e.Id = tsv[0]
	e.Class = tsv[1]
	e.Analysis = strings.Split(tsv[2], ",")
	e.Diagnosis = strings.Split(tsv[3], ",")
	return
}

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

func tsvToEntities(path string) (entities []Entity, err error) {
	tsv, err := readTsv(path)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not read file %s", path))
		return
	}
	for _, row := range tsv {
		var entity Entity
		if err = entity.tsvToEntity(row); err != nil {
			return
		}
		entities = append(entities, entity)
	}
	return
}
