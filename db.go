package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func (s Session) getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", s.DbHost, s.DbPort, s.DbUser, s.DbPassword, s.DbName)
}

func (s Session) dbExec(query string) (err error) {
	var stmt *sql.Stmt
	stmt, err = s.DbConnection.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	return
}

func (s Session) dbQuery(query string) (rows *sql.Rows, err error) {
	var stmt *sql.Stmt
	stmt, err = s.DbConnection.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err = stmt.Query()
	return
}

func (s Session) checkTableExists(table string) (err error) {
	var rows *sql.Rows
	if rows, err = s.dbQuery(fmt.Sprintf("SELECT * FROM %s", table)); err != nil {
		return
	}
	defer rows.Close()
	return
}

func (s Session) createTable(table string) (err error) {
	if err = s.dbExec(fmt.Sprintf("CREATE TABLE %s (id varchar(20) NOT NULL, ensembl varchar(20) NOT NULL, class varchar(10) NOT NULL, %s boolean, PRIMARY KEY (id))", table, strings.Join(analyses, " boolean, "))); err != nil {
		return
	}
	log.Println(fmt.Sprintf("Table %s was added to database %s.", table, s.DbName))
	return
}

func (s *Session) initDb() (err error) {
	if s.DbConnection, err = sql.Open("postgres", s.getConnectionString()); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not establish connection to database %s.", s.DbName))
		return
	}
	if err = s.DbConnection.Ping(); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not establish connection to database %s.", s.DbName))
	}
	for _, list := range lists {
		if err = s.checkTableExists(fmt.Sprintf("list_%s", list)); err == nil {
			continue
		} else {
			if err = s.createTable(fmt.Sprintf("list_%s", list)); err != nil {
				err = errors.Wrap(err, fmt.Sprintf("Could not create table %s", list))
				return
			}
		}
	}
	return
}

func (s Session) checkEntityExists(table string, entity Entity) (exists bool) {
	rows, err := s.dbQuery(fmt.Sprintf("SELECT EXISTS (SELECT id FROM %s WHERE id = '%s');", table, entity.Id))
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

func (s Session) updateEntity(table string, entity Entity) (err error) {
	var analyses []string
	for key, _ := range entity.Analyses {
		analyses = append(analyses, key)
	}
	updateString := strings.Join(analyses, " = true, ")
	err = s.dbExec(fmt.Sprintf("UPDATE %s SET %s = true WHERE id = '%s';", table, updateString, entity.Id))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not update entity %s in table %s", entity.Id, table))
		return
	}
	log.Printf("Entity %s in table %s was updated.", entity.Id, table)
	return
}

func (s Session) addEntity(table string, entity Entity) (err error) {
	err = s.dbExec(fmt.Sprintf("INSERT INTO %s (id, ensembl, class, snv, cnv, sv, pindel) VALUES ('%s', '%s', '%s', %t, %t, %t, %t)", table, entity.Id, entity.EnsemblId, entity.Class, entity.Analyses["snv"], entity.Analyses["cnv"], entity.Analyses["sv"], entity.Analyses["pindel"]))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not add entity %s to table %s", entity.Id, table))
		return
	}
	log.Printf("Entity %s was added to table %s", entity.Id, table)
	return
}

func (s Session) getEntityList(table string, column string) (ids []Entity, err error) {
	log.Printf("Retriewing %s gene list from %s", column, table)
	var rows *sql.Rows
	if table == "list_aml_ext" {
		rows, err = s.dbQuery(fmt.Sprintf("SELECT id, ensembl FROM %s FULL OUTER JOIN list_aml USING (id, ensembl) WHERE %s.%s = true OR list_aml.%s = true;", table, table, column, column))
	} else {
		rows, err = s.dbQuery(fmt.Sprintf("SELECT id, ensembl FROM %s WHERE %s = true;", table, column))
	}
	defer rows.Close()
	var id Entity
	for rows.Next() {
		err = rows.Scan(&id.Id, &id.EnsemblId)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}
