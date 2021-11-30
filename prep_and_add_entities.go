package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func setEntityAnalysesOrLists(slice []string) (field map[string]bool) {
	field = make(map[string]bool)
	for _, element := range slice {
		field[element] = true
	}
	return
}

func tsvToEntity(tsv []string) (entity Entity, err error) {
	if len(tsv) != 5 {
		err = errors.New("Input tsv needs to have 5 columns")
		return
	}
	entity.Id = tsv[0]
	entity.Class = tsv[1]
	entity.Analyses = setEntityAnalysesOrLists(strings.Split(tsv[2], ","))
	entity.Lists = setEntityAnalysesOrLists(strings.Split(tsv[3], ","))
	if tsv[4] == "all" || tsv[4] == "All" {
		entity.AllFusions = true
	}
	return
}

func (s Session) checkAndAddEntity(list string, entity Entity) (err error) {
	if s.checkEntityExists(fmt.Sprintf("list_%s", list), entity) {
		if err = s.updateEntity(fmt.Sprintf("list_%s", list), entity); err != nil {
			return
		}
	} else {
		if entity.Ensembl38Id, err = s.getEnsemblId(entity.Id, "38"); err != nil {
			return
		}
		if entity.Ensembl37Id, err = s.getEnsemblId(entity.Id, "37"); err != nil {
			return
		}
		if entity.Ensembl38Id != "" || entity.Ensembl37Id != "" {
			err = s.addEntity(fmt.Sprintf("list_%s", list), entity)
			return
		} else if unknownIds[entity.Id] != "" {
			entity.Ensembl38Id = unknownIds[entity.Id]
			err = s.addEntity(fmt.Sprintf("list_%s", list), entity)
			return
		} else {
			log.Printf("Did not find Ensembl entry for %s", entity.Id)
		}
	}
	return
}

func (s Session) prepAndAddFusions(list string, entity Entity) (err error) {
	if entity.Analyses["sv"] && entity.AllFusions {
		log.Printf("Starting to scrape all fusion partners for %s", entity.Id)
		var genes []string
		genes, err = s.getFusionGenes(entity)
		if err != nil {
			return
		}
		for _, gene := range genes {
			fusionEntity := Entity{
				Id:       gene,
				Class:    "Gene",
				Analyses: map[string]bool{"sv": true},
			}
			if err = s.checkAndAddEntity(list, fusionEntity); err != nil {
				return
			}
		}
	}
	return
}

func (s Session) prepAndAddEntities() (err error) {
	tsv, err := readTsv(s.Tsv)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not read file %s", s.Tsv))
		return
	}
	for _, row := range tsv {
		var entity Entity
		if entity, err = tsvToEntity(row); err != nil {
			return
		}
		for list, _ := range entity.Lists {
			if err = s.checkAndAddEntity(list, entity); err != nil {
				return
			}
			if err = s.prepAndAddFusions(list, entity); err != nil {
				return
			}
		}
	}
	return
}
