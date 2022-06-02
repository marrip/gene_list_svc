package cmd

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
)

func getMockDb(route string) {
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
	case "cannotCreateNewRow":
		prep := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO "existing_table" (id, ensembl_id_38, ensembl_id_37, class, chromosome, start, "end", cnv, pindel, snv, sv) VALUES ('GENE1', '', '', 'gene', '', '', '', true, false, true, false)`))
		prep.ExpectExec().WithArgs().WillReturnError(fmt.Errorf("Something went wrong"))
	case "cannotCreateNewTable":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "nonexistent_table"`)).WillReturnError(fmt.Errorf("Something went wrong"))
	case "cannotGetRegions":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, ensembl_id_38, ensembl_id_37, class, chromosome, start, "end" FROM "test" WHERE snv = true;`)).WillReturnError(fmt.Errorf("Something went wrong"))
	case "cannotGetTables":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`)).WillReturnError(fmt.Errorf("Something went wrong"))
	case "checkAndCreateNewTable":
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "new_table"`)).WillReturnError(fmt.Errorf("Something went wrong"))
		prep := mock.ExpectPrepare(regexp.QuoteMeta(`CREATE TABLE "new_table" (id varchar(20) NOT NULL, ensembl_id_38 varchar(20) NOT NULL, ensembl_id_37 varchar(20) NOT NULL, class varchar(10) NOT NULL, chromosome varchar(2) NOT NULL, start varchar(10) NOT NULL, "end" varchar(10) NOT NULL, cnv boolean, pindel boolean, snv boolean, sv boolean, PRIMARY KEY (id))`))
		prep.ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
	case "createNewRow":
		prep := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO "existing_table" (id, ensembl_id_38, ensembl_id_37, class, chromosome, start, "end", cnv, pindel, snv, sv) VALUES ('GENE1', '', '', 'gene', '', '', '', true, false, true, false)`))
		prep.ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	case "createNewTable":
		prep := mock.ExpectPrepare(regexp.QuoteMeta(`CREATE TABLE "new_table" (id varchar(20) NOT NULL, ensembl_id_38 varchar(20) NOT NULL, ensembl_id_37 varchar(20) NOT NULL, class varchar(10) NOT NULL, chromosome varchar(2) NOT NULL, start varchar(10) NOT NULL, "end" varchar(10) NOT NULL, cnv boolean, pindel boolean, snv boolean, sv boolean, PRIMARY KEY (id))`))
		prep.ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
	case "default":
	case "getRegions":
		rows := sqlmock.NewRows([]string{"id", "ensembl_id_38", "ensembl_id_37", "class", "chromosome", "start", "end"}).AddRow("GENE1", "ENSG001", "ENSG001", "gene", "1", "1", "100")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, ensembl_id_38, ensembl_id_37, class, chromosome, start, "end" FROM "test" WHERE snv = true;`)).WillReturnRows(rows)
	case "getTables":
		rows := sqlmock.NewRows([]string{"table_name"}).AddRow("test")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`)).WillReturnRows(rows)
	case "regionExists":
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (SELECT id FROM "existing_table" WHERE id = 'GENE1');`)).WillReturnRows(rows)
	case "tableExists":
		rows := sqlmock.NewRows([]string{"id", "ensembl_id_38", "ensembl_id_37", "class", "chromosome", "start", "end", "snv", "cnv", "sv", "pindel"})
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "existing_table"`)).WillReturnRows(rows)
	case "updateRow":
		prep := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE "existing_table" SET cnv = true WHERE id = 'GENE1';`))
		prep.ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	}
	return
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

func TestCheckRegionExists(t *testing.T) {
	var cases = map[string]struct {
		table  string
		region DbTableRow
		route  string
		result bool
	}{
		"Region exists": {
			"existing_table",
			DbTableRow{
				Id: "GENE1",
			},
			"regionExists",
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			result := session.Db.Connection.checkRegionExists(c.table, c.region)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestUpdateRow(t *testing.T) {
	var cases = map[string]struct {
		table   string
		region  DbTableRow
		route   string
		wantErr bool
	}{
		"Update row successfully": {
			"existing_table",
			DbTableRow{
				Analyses: map[string]struct{}{
					"cnv": struct{}{},
				},
				Id: "GENE1",
			},
			"updateRow",
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			err := session.Db.Connection.updateRow(c.table, c.region)
			checkError(t, err, c.wantErr)
		})
	}
}

func TestAddNewRow(t *testing.T) {
	var cases = map[string]struct {
		table   string
		region  DbTableRow
		route   string
		wantErr bool
	}{
		"Add row successfully": {
			"existing_table",
			DbTableRow{
				Analyses: map[string]struct{}{
					"cnv": struct{}{},
					"snv": struct{}{},
				},
				Class: "gene",
				Id:    "GENE1",
			},
			"createNewRow",
			false,
		},
		"Row could not be added": {
			"existing_table",
			DbTableRow{
				Analyses: map[string]struct{}{
					"cnv": struct{}{},
					"snv": struct{}{},
				},
				Class: "gene",
				Id:    "GENE1",
			},
			"cannotCreateNewRow",
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			err := session.Db.Connection.addNewRow(c.table, c.region)
			checkError(t, err, c.wantErr)
		})
	}
}

func TestGetAnalysis(t *testing.T) {
	var cases = map[string]struct {
		analysis string
		result   bool
	}{
		"Return true": {
			"cnv",
			true,
		},
		"Return false": {
			"pindel",
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			row := DbTableRow{
				Analyses: map[string]struct{}{
					"cnv": struct{}{},
				},
			}
			result := row.getAnalysis(c.analysis)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetTables(t *testing.T) {
	var cases = map[string]struct {
		route   string
		result  map[string]struct{}
		wantErr bool
	}{
		"Get tables successfully": {
			"getTables",
			map[string]struct{}{
				"test": {},
			},
			false,
		},
		"Could not get tables": {
			"cannotGetTables",
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			result, err := session.Db.Connection.getTables()
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetRegions(t *testing.T) {
	var cases = map[string]struct {
		route   string
		result  []DbTableRow
		wantErr bool
	}{
		"Get regions successfully": {
			"getRegions",
			[]DbTableRow{
				DbTableRow{
					Id:          "GENE1",
					EnsemblId37: "ENSG001",
					EnsemblId38: "ENSG001",
					Class:       "gene",
					Chromosome:  "1",
					Start:       "1",
					End:         "100",
				},
			},
			false,
		},
		"Could not get regions": {
			"cannotGetRegions",
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			getMockDb(c.route)
			session.Analysis = "snv"
			session.Tables = []string{"test"}
			result, err := session.Db.Connection.getRegions()
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
