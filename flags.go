package main

import (
	"flag"
	"strings"
)

func (s *Session) readFlags() (err error) {
	var selectors string
	flag.StringVar(&s.Path, "tsv", "", "Path to tsv file containing gene list.")
	flag.StringVar(&selectors, "select", "", "Comma separated list of selectors for gene list.")
	flag.Parse()
	s.Selectors = strings.Split(selectors, ",")
	return
}
