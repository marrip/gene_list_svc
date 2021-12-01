package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func sliceContainsString(arg string, slice []string) bool {
	for _, element := range slice {
		if arg == element {
			return true
		}
	}
	return false
}

func (s Session) checkListAndAnalysis() (err error) {
	if s.Analysis == "" && s.List == "" {
		return
	} else if (s.Analysis == "" && s.List != "") || (s.Analysis != "" && s.List == "") {
		err = errors.New("The flags -list and -analysis need to be used together")
	} else if !sliceContainsString(s.List, lists) {
		err = errors.New(fmt.Sprintf("Please choose a valid list (%s)", strings.Join(lists, ", ")))
	} else if !sliceContainsString(s.Analysis, analyses) {
		err = errors.New(fmt.Sprintf("Please choose a valid analysis (%s)", strings.Join(analyses, ", ")))
	}
	return
}

func (s *Session) checkBuild() (err error) {
	if s.Build == "" {
		s.Build = "38"
		return
	} else if sliceContainsString(s.Build, builds) {
		return
	} else {
		err = errors.New(fmt.Sprintf("Please choose a valid genome build (%s)", strings.Join(builds, ", ")))
	}
	return
}

func (s *Session) modifyFlagInput() {
	if s.Analysis != "" && s.List != "" {
		if s.Bed == "" {
			s.Bed = fmt.Sprintf("%s_%s_%s.bed", s.List, s.Analysis, s.Build)
		}
		if s.List != "" {
			s.DbList = fmt.Sprintf("list_%s", s.List)
		}
	}
}

func (s *Session) readFlags() (err error) {
	flag.StringVar(&s.Analysis, "analysis", "", fmt.Sprintf("Select analysis (%s).", strings.Join(analyses, ", ")))
	flag.StringVar(&s.Bed, "bed", "", "Output bed file name (default: [list]_[analysis]_[build].bed).")
	flag.StringVar(&s.Build, "build", "", fmt.Sprintf("Select genome build (%s; default: 38).", strings.Join(builds, ", ")))
	flag.StringVar(&s.List, "list", "", fmt.Sprintf("Select gene list (%s).", strings.Join(lists, ", ")))
	flag.StringVar(&s.Tsv, "tsv", "", "Path to tsv file containing gene list.")
	flag.Parse()
	if err = s.checkListAndAnalysis(); err != nil {
		return
	}
	if err = s.checkBuild(); err != nil {
		return
	}
	s.modifyFlagInput()
	return
}

func initSession() (session Session, err error) {
	if err = env.Parse(&session); err != nil {
		err = errors.Wrap(err, "Could not read env.")
	}
	log.Println("Successfully read env.")
	if err = session.readFlags(); err != nil {
		return
	}
	log.Println("Successfully read flags.")
	if err = session.initDb(); err != nil {
		return
	}
	log.Println("Successfully initialized database.")
	return
}
