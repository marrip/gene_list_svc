package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

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
	} else if err = session.Db.Connection.createTable(table); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not create table %s", table))
		return
	}
	return
}

func (d dbConnection) checkTableExists(table string) (exists bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return
	}
	defer rows.Close()
	exists = true
	return
}

func (d dbConnection) createTable(table string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := fmt.Sprintf("CREATE TABLE ? (id varchar(20) NOT NULL, ensembl_id_38 varchar(20) NOT NULL, ensembl_id_37 varchar(20) NOT NULL, class varchar(10) NOT NULL, chromosome varchar(2) NOT NULL, start varchar(10) NOT NULL, end varchar(10) NOT NULL, %s boolean, PRIMARY KEY (id))", strings.Join(getAnalyses(analyses), " boolean, "))
	var stmt *sql.Stmt
	stmt, err = d.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, table)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Table %s was added to database %s.", table, session.Db.Name))
	return
}

func getAnalyses(analyses map[string]struct{}) (analysesSlice []string) {
	for analysis, _ := range analyses {
		analysesSlice = append(analysesSlice, analysis)
	}
	sort.Strings(analysesSlice)
	return
}

func (d dbConnection) checkRegionExists(table string, region DbTableRow) (exists bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("SELECT EXISTS (SELECT id FROM %s WHERE id = '%s');", table, region.Id))
	if err != nil {
		fmt.Printf("%v", err)
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

func (d dbConnection) updateRow(table string, region DbTableRow) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := fmt.Sprintf("UPDATE ? SET %s = true WHERE id = '?';", strings.Join(getAnalyses(region.Analyses), " = true, "))
	var stmt *sql.Stmt
	stmt, err = d.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, table, region.Id)
	if err != nil {
		return
	}
	log.Printf("Region %s in table %s was updated.", region.Id, table)
	return
}

func (d dbConnection) addNewRow(table string, region DbTableRow) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO ? (id, ensembl_id_38, ensembl_id_37, class, chromosome, start, end, cnv, pindel, snv, sv) VALUES ('?', '?', '?', '?', '?', '?', '?', ?, ?, ?, ?)"
	var stmt *sql.Stmt
	stmt, err = d.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, table, region.Id, region.EnsemblId38, region.EnsemblId37, region.Class, region.Chromosome, region.Start, region.End, region.getAnalysis("cnv"), region.getAnalysis("pindel"), region.getAnalysis("snv"), region.getAnalysis("sv"))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not add region %s to table %s", region.Id, table))
		return
	}
	log.Printf("Region %s was added to table %s.", region.Id, table)
	return
}

func (d DbTableRow) getAnalysis(analysis string) (include bool) {
	_, include = d.Analyses[analysis]
	return
}
