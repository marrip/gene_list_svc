package main

import (
	"database/sql"
)

type Session struct {
	DbHost       string `env:"DB_HOST" envDefault:"localhost"`
	DbPort       int    `env:"DB_PORT" envDefault:"5432"`
	DbUser       string `env:"DB_USER"`
	DbPassword   string `env:"DB_PASSWORD"`
	DbName       string `env:"DB_NAME" envDefault:"gene_list"`
	DbConnection *sql.DB
	Path         string
}

type Entity struct {
	Id        string
	Class     string
	Analysis  []string
	Diagnosis []string
}

type Partner struct {
	First  string
	Second string
}
