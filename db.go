package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/pkg/errors"
)

func fillSlice(values []string, valueMap map[string]bool) (boolSlice []bool) {
	for _, value := range values {
		valueMap[value] = true
	}
	var keys []string
	for key, _ := range valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		boolSlice = append(boolSlice, valueMap[key])
	}
	return
}

func (e Entity) getAnalysisBool() (analysisBoolsSlice []bool) {
	analysisBoolsMap := map[string]bool{
		"cnv":    false,
		"pindel": false,
		"snv":    false,
		"sv":     false,
	}
	analysisBoolsSlice = fillSlice(e.Analysis, analysisBoolsMap)
	return
}

func (e Entity) getDiagnosisBool() (diagnosisBoolsSlice []bool) {
	diagnosisBoolsMap := map[string]bool{
		"list_all":     false,
		"list_aml":     false,
		"list_aml_ext": false,
	}
	diagnosisBoolsSlice = fillSlice(e.Diagnosis, diagnosisBoolsMap)
	return
}

func (s Session) addEntity(entity Entity) (err error) {
	var stmt *sql.Stmt
	analysisSlice := entity.getAnalysisBool()
	diagnosisSlice := entity.getDiagnosisBool()
	stmt, err = s.DbConnection.Prepare(fmt.Sprintf("INSERT INTO entity (id, class, cnv, pindel, snv, sv, list_all, list_aml, list_aml_ext) VALUES ('%s', '%s', %t, %t, %t, %t, %t, %t, %t)", entity.Id, entity.Class, analysisSlice[0], analysisSlice[1], analysisSlice[2], analysisSlice[3], diagnosisSlice[0], diagnosisSlice[1], diagnosisSlice[2]))
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not add entity %s", entity.Id))
	}
	return
}

func (s Session) addEntities(entities []Entity) {
	for _, entity := range entities {
		if err := s.addEntity(entity); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("Could not add entity %s", entity.Id))
			log.Printf("%v", err)
		}
	}
	return
}
