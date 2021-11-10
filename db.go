package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func (s Session) checkEntityExists(entity Entity) (exists bool) {
	stmt, err := s.DbConnection.Prepare(fmt.Sprintf("SELECT EXISTS (SELECT id FROM entity WHERE id = '%s');", entity.Id))
	if err != nil {
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	rows, err = stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return
		}
	}
	return
}

func (s Session) updateEntity(entity Entity) (err error) {
	updateString := strings.Join(append(entity.Analysis, entity.Diagnosis...), " = true, ")
	stmt, err := s.DbConnection.Prepare(fmt.Sprintf("UPDATE entity SET %s = true WHERE id = '%s';", updateString, entity.Id))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not update entity %s", entity.Id))
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not update entity %s", entity.Id))
		return
	}
	log.Println(fmt.Sprintf("Entity %s was updated.", entity.Id))
	return
}

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
		return
	}
	log.Println(fmt.Sprintf("Entity %s was added to database.", entity.Id))
	return
}

func (s Session) addEntities(entities []Entity) {
	for _, entity := range entities {
		if s.checkEntityExists(entity) {
			err := s.updateEntity(entity)
			if err != nil {
				log.Printf("%v", err)
			}
		} else {
			if err := s.addEntity(entity); err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Could not add entity %s", entity.Id))
				log.Printf("%v", err)
			}
		}
	}
	return
}

func (s Session) getEntityList(columns []string) (ids []string, err error) {
	var stmt *sql.Stmt
	queryString := strings.Join(columns, " = true AND ")
	stmt, err = s.DbConnection.Prepare(fmt.Sprintf("SELECT id FROM entity WHERE %s = true;", queryString))
	if err != nil {
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	rows, err = stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

func (s Session) setGetDb() (err error) {
	if s.Path != "" {
		var entities []Entity
		entities, err = tsvToEntities(s.Path)
		if err != nil {
			return
		}
		s.addEntities(entities)
	}
	if len(s.Selectors) > 0 {
		var ids []string
		ids, err = s.getEntityList(s.Selectors)
		if err != nil {
			return
		}
		fmt.Println(ids)
	}
	return
}
