package main

import (
	"log"
)

func main() {
	s, err := initSession()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err = s.initDb(); err != nil {
		log.Fatalf("%v", err)
	}
	defer s.DbConnection.Close()
	//entities, err := tsvToEntities("test.tsv")
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}
	//s.addEntities(entities)
	//log.Println("All done")
}
