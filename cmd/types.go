package cmd

import (
	"database/sql"
)

type Session struct {
	Db       database
	Web      web
	Analysis string
	Bed      string
	Build    string
	Chr      bool
	Tables   []string
	Tsv      string
}

type database struct {
	Connection DbConnection
	Host       string `env:"DB_HOST" envDefault:"localhost"`
	Name       string `env:"DB_NAME" envDefault:"gene_list"`
	Password   string `env:"DB_PASSWORD"`
	Port       int    `env:"DB_PORT" envDefault:"5432"`
	User       string `env:"DB_USER"`
}

type DbConnection interface {
	addNewRow(table string, region DbTableRow) (err error)
	checkRegionExists(table string, region DbTableRow) (exists bool)
	checkTableExists(table string) (exists bool)
	createTable(table string) (err error)
	getRegions() (regions []DbTableRow, err error)
	getTables() (tables map[string]struct{}, err error)
	updateRow(table string, region DbTableRow) (err error)
}

type dbConnection struct {
	db *sql.DB
}

type web struct {
	AtlasGO   string `env:"ATLAS_ROOT_URL" envDefault:"http://atlasgeneticsoncology.org"`
	Ensembl38 string `env:"ENSEMBL_38_REST_URL" envDefault:"https://rest.ensembl.org"`
	Ensembl37 string `env:"ENSEMBL_37_REST_URL" envDefault:"https://grch37.rest.ensembl.org"`
}

type DbTableRow struct {
	Analyses        map[string]struct{}
	End             string
	EnsemblId38     string
	EnsemblId37     string
	Chromosome      string
	Class           string
	Id              string
	IncludePartners bool
	Start           string
	Tables          []string
}

type EnsemblGeneObj struct {
	Chromosome  string            `json:"seq_region_name"`
	EnsemblId   string            `json:"id"`
	Id          string            `json:"display_name"`
	Transcripts []EnsemblTransObj `json:"Transcript"`
}

type EnsemblTransObj struct {
	EnsemblId string           `json:"id"`
	Exons     []EnsemblBaseObj `json:"Exon"`
}

type EnsemblBaseObj struct {
	Annotation string
	Chromosome string `json:"seq_region_name"`
	End        int    `json:"end"`
	EnsemblId  string `json:"id"`
	Start      int    `json:"start"`
	Transcript string
}
