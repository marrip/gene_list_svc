package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func (s *Session) initDbConnection() (err error) {
	var connection *sql.DB
	if connection, err = sql.Open("postgres", getConnectionString()); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not establish connection to database %s.", s.Db.Name))
		return
	}
	if err = connection.Ping(); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not establish connection to database %s.", s.Db.Name))

	}
	s.Db.Connection = dbConnection{connection}
	return
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", session.Db.Host, session.Db.Port, session.Db.User, session.Db.Password, session.Db.Name)
}

func ensureTableExists(table string) (err error) {
	if session.Db.Connection.checkTableExists(table) {
		return
	} else if err = createTable(table); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not create table %s", table))
		return
	}
	return
}

func (d dbConnection) checkTableExists(table string) (exists bool) {
	stmt, err := d.db.Prepare(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	exists = true
	return
}

func createTable(table string) (err error) {
	if err = session.Db.Connection.dbExec(fmt.Sprintf("CREATE TABLE %s (id varchar(20) NOT NULL, ensemble_id_38 varchar(20) NOT NULL, ensembl_id_37 varchar(20) NOT NULL, class varchar(10) NOT NULL, chromosome varchar(2) NOT NULL, start varchar(10) NOT NULL, end varchar(10) NOT NULL,  %s boolean, PRIMARY KEY (id))", table, strings.Join(getAnalyses(), " boolean, "))); err != nil {
		return
	}
	log.Println(fmt.Sprintf("Table %s was added to database %s.", table, session.Db.Name))
	return
}

func (d dbConnection) dbExec(query string) (err error) {
	var stmt *sql.Stmt
	stmt, err = d.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	return
}

func getAnalyses() (analysesSlice []string) {
	for analysis, _ := range analyses {
		analysesSlice = append(analysesSlice, analysis)
	}
	return
}
