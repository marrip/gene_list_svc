package main

import (
	"database/sql"
)

type Session struct {
	DbHost           string `env:"DB_HOST" envDefault:"localhost"`
	DbPort           int    `env:"DB_PORT" envDefault:"5432"`
	DbUser           string `env:"DB_USER"`
	DbPassword       string `env:"DB_PASSWORD"`
	DbName           string `env:"DB_NAME" envDefault:"gene_list"`
	DbList           string
	DbConnection     *sql.DB
	AtlasRootUrl     string `env:"ATLAS_ROOT_URL" envDefault:"http://atlasgeneticsoncology.org"`
	Ensembl38RestUrl string `env:"ENSEMBL_38_REST_URL" envDefault:"https://rest.ensembl.org"`
	Ensembl37RestUrl string `env:"ENSEMBL_37_REST_URL" envDefault:"https://grch37.rest.ensembl.org"`
	Analysis         string
	Bed              string
	Build            string
	List             string
	Tsv              string
}

type Entity struct {
	Id          string
	Ensembl38Id string
	Ensembl37Id string
	Class       string
	Analyses    map[string]bool
	AllFusions  bool
	Lists       map[string]bool
	Diagnosis   []string
}

type Region struct {
	Gene       string
	Id         string
	Chromosome string
	Start      int
	End        int
}

// Ensembl Json Objects

type EnsemblGeneObj struct {
	Id          string            `json:"display_name"`
	EnsemblId   string            `json:"id"`
	Type        string            `json:"type"`
	Chromosome  string            `json:"seq_region_name"`
	Start       int               `json:"start"`
	End         int               `json:"end"`
	Transcripts []EnsemblTransObj `json:"Transcript"`
}

type EnsemblTransObj struct {
	EnsemblId  string           `json:"id"`
	Chromosome string           `json:"seq_region_name"`
	Start      int              `json:"start"`
	End        int              `json:"end"`
	Exons      []EnsemblExonObj `json:"Exon"`
}

type EnsemblExonObj struct {
	EnsemblId  string `json:"id"`
	Chromosome string `json:"seq_region_name"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
}
