package main

import (
	"flag"

	"github.com/pkg/errors"
)

func (s *Session) readFlags() (err error) {
	flag.StringVar(&s.Path, "tsv", "", "Path to tsv file containing gene list.")
	flag.Parse()
	if s.Path == "" {
		err = errors.New("Please indicate a tsv file as input")
	}
	return
}
