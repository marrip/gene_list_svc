package main

import (
	"log"
)

func main() {
	s, err := initSession()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer s.DbConnection.Close()
	if s.Tsv != "" {
		err = s.prepAndAddEntities()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	if s.List != "" && s.Analysis != "" {
		err = s.getCoordinatesforEntities()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
