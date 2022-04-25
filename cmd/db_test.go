package cmd

import (
	//"database/sql/driver"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
)

func getMockDb(route string) sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Could not create mock database. Got Error\n%v", err))
	}
	session = Session{
		Db: database{
			Connection: dbConnection{
				db: db,
			},
			Host:     "localhost",
			Name:     "my_list",
			Password: "sodamnsecret",
			Port:     5432,
			User:     "itsme",
		},
	}
	switch route {
	case "cannotCreateNewTable":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM nonexistent_table`)).WillReturnError(fmt.Errorf("Something went wrong"))
	case "checkAndCreateNewTable":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM new_table`)).WillReturnError(fmt.Errorf("Something went wrong"))
		prep := mock.ExpectPrepare("CREATE TABLE \\? \\(id varchar\\(20\\) NOT NULL, ensembl_id_38 varchar\\(20\\) NOT NULL, ensembl_id_37 varchar\\(20\\) NOT NULL, class varchar\\(10\\) NOT NULL, chromosome varchar\\(2\\) NOT NULL, start varchar\\(10\\) NOT NULL, end varchar\\(10\\) NOT NULL, cnv boolean, pindel boolean, snv boolean, sv boolean, PRIMARY KEY \\(id\\)\\)")
		prep.ExpectExec().WithArgs("new_table").WillReturnResult(sqlmock.NewResult(0, 0))
	case "createNewTable":
		prep := mock.ExpectPrepare("CREATE TABLE \\? \\(id varchar\\(20\\) NOT NULL, ensembl_id_38 varchar\\(20\\) NOT NULL, ensembl_id_37 varchar\\(20\\) NOT NULL, class varchar\\(10\\) NOT NULL, chromosome varchar\\(2\\) NOT NULL, start varchar\\(10\\) NOT NULL, end varchar\\(10\\) NOT NULL, cnv boolean, pindel boolean, snv boolean, sv boolean, PRIMARY KEY \\(id\\)\\)")
		prep.ExpectExec().WithArgs("new_table").WillReturnResult(sqlmock.NewResult(0, 0))
	case "default":
	case "tableExists":
		rows := sqlmock.NewRows([]string{"id", "ensembl_id_38", "ensembl_id_37", "class", "chromosome", "start", "end", "snv", "cnv", "sv", "pindel"})
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM existing_table`)).WillReturnRows(rows)
	}
	return mock
}

func TestGetConnectionString(t *testing.T) {
	var cases = map[string]struct {
		result string
	}{
		"Return connection string": {
			"host=localhost port=5432 user=itsme password=sodamnsecret dbname=my_list sslmode=disable",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb("default")
			result := getConnectionString()
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestEnsureTableExists(t *testing.T) {
	var cases = map[string]struct {
		table   string
		route   string
		wantErr bool
	}{
		"Table exists": {
			"existing_table",
			"tableExists",
			false,
		},
		"Table does not exist": {
			"new_table",
			"checkAndCreateNewTable",
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			err := ensureTableExists(c.table)
			checkError(t, err, c.wantErr)
		})
	}
}

func TestCheckTableExists(t *testing.T) {
	var cases = map[string]struct {
		table  string
		route  string
		result bool
	}{
		"Table exists": {
			"existing_table",
			"tableExists",
			true,
		},
		"Table does not exist": {
			"new_table",
			"createNewTable",
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			result := session.Db.Connection.checkTableExists(c.table)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestCreateTable(t *testing.T) {
	var cases = map[string]struct {
		table   string
		route   string
		wantErr bool
	}{
		"Add table successfully": {
			"new_table",
			"createNewTable",
			false,
		},
		"Table could not be added": {
			"nonexistent_table",
			"cannotCreateNewTable",
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			err := session.Db.Connection.createTable(c.table)
			checkError(t, err, c.wantErr)
		})
	}
}

func TestGetAnalyses(t *testing.T) {
	var cases = map[string]struct {
		analyses map[string]struct{}
		result   []string
	}{
		"Return all existing analyses": {
			analyses,
			[]string{
				"cnv",
				"pindel",
				"snv",
				"sv",
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := getAnalyses(c.analyses)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
