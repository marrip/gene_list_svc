package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func initSession() (session Session, err error) {
	if err = env.Parse(&session); err != nil {
		err = errors.Wrap(err, "Could not read env.")
	}
	log.Println("Successfully read env.")
	if err = session.readFlags(); err != nil {
		return
	}
	log.Println("Successfully read flags.")
	return
}

func (s Session) getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", s.DbHost, s.DbPort, s.DbUser, s.DbPassword, s.DbName)
}

func (s Session) checkTableExists(table string) (err error) {
	var stmt *sql.Stmt
	stmt, err = s.DbConnection.Prepare(fmt.Sprintf("SELECT * FROM %s", table)) // Does not work!!!
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
	return
}

func (s Session) createDbTable() (err error) {
	for key, cmd := range tables {
		if err = s.checkTableExists(key); err == nil {
			return
		} else {
			var stmt *sql.Stmt
			stmt, err = s.DbConnection.Prepare(fmt.Sprintf("CREATE TABLE %s (%s)", key, cmd))
			if err != nil {
				return
			}
			defer stmt.Close()
			_, err = stmt.Exec()
			if err != nil {
				return
			}
			log.Println(fmt.Sprintf("Table %s was added to database %s.", key, s.DbName))
		}
	}
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
	err = s.createDbTable()
	return
}
