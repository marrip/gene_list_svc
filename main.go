package main

import (
	//"log"
	"github.com/marrip/gene_list_svc/cmd"
)

func main() {
	cmd.Execute()
	//s, err := initSession()
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}
	//defer s.DbConnection.Close()
	//if s.Tsv != "" {
	//	err = s.prepAndAddEntities()
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//}
	//if s.Table != "" && s.Analysis != "" {
	//	err = s.getCoordinatesforEntities()
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//}
}
