package main

import (
	"database/sql"
)

type Session struct {
	DbHost         string `env:"DB_HOST" envDefault:"localhost"`
	DbPort         int    `env:"DB_PORT" envDefault:"5432"`
	DbUser         string `env:"DB_USER"`
	DbPassword     string `env:"DB_PASSWORD"`
	DbName         string `env:"DB_NAME" envDefault:"gene_list"`
	DbList         string
	DbConnection   *sql.DB
	AtlasRootUrl   string `env:"ATLAS_ROOT_URL" envDefault:"http://atlasgeneticsoncology.org"`
	EnsemblRestUrl string `env:"ENSEMBL_REST_URL" envDefault:"https://rest.ensembl.org"`
	Analysis       string
	Bed            string
	List           string
	Tsv            string
}

type Entity struct {
	Id         string
	EnsemblId  string
	Class      string
	Analyses   map[string]bool
	AllFusions bool
	Lists      map[string]bool
	Diagnosis  []string
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
	Exons       []EnsemblExonObj
}

type EnsemblTransObj struct {
	EnsemblId  string           `json:"id"`
	Chromosome string           `json:"seq_region_name"`
	Start      int              `json:"start"`
	End        int              `json:"end"`
	Exons      []EnsemblExonObj `json:"Exon"`
}

type EnsemblExonObj struct {
	EnsemblId    string `json:"id"`
	TranscriptId string
	Chromosome   string `json:"seq_region_name"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
}
