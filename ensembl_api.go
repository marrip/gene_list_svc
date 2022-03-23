package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

func (s Session) getBuildUrl() string {
	if s.Build == "37" {
		return s.Ensembl37RestUrl
	}
	return s.Ensembl38RestUrl
}

func (s Session) getEnsemblIdUrl(gene string) string {
	return fmt.Sprintf("%s/xrefs/symbol/homo_sapiens/%s?content-type=application/json", s.getBuildUrl(), gene)
}

func (s Session) getEnsemblSmallGeneUrl(id string) string {
	return fmt.Sprintf("%s/lookup/id/%s?content-type=application/json", s.getBuildUrl(), id)
}

func (s Session) checkEnsemblIdChromosome(id string, build string) (isMain bool, err error) {
	body, err := sendHttpRequest(s.getEnsemblSmallGeneUrl(id))
	if err != nil {
		return
	}
	var jsonObj EnsemblGeneObj
	json.Unmarshal(body, &jsonObj)
	for _, chromosome := range chromosomes {
		if chromosome == jsonObj.Chromosome {
			isMain = true
		}
	}
	return
}

func (s Session) getEnsemblId(gene string, build string) (id string, err error) {
	ensemblRegex := regexp.MustCompile("ENSG")
	body, err := sendHttpRequest(s.getEnsemblIdUrl(gene))
	if err != nil {
		return
	}
	var jsonObj []EnsemblGeneObj
	json.Unmarshal(body, &jsonObj)
	for _, element := range jsonObj {
		if ensemblRegex.MatchString(element.EnsemblId) {
			var isMain bool
			isMain, err = s.checkEnsemblIdChromosome(element.EnsemblId, build)
			if err != nil {
				return
			}
			if isMain {
				id = element.EnsemblId
				return
			}
		}
	}
	return
}

func (s Session) getEnsemblIds(entities []Entity) (updates []Entity, err error) {
	for _, entity := range entities {
		entity.Ensembl38Id, err = s.getEnsemblId(entity.Id, "38")
		if err != nil {
			return
		}
		entity.Ensembl37Id, err = s.getEnsemblId(entity.Id, "37")
		if err != nil {
			return
		}
		if entity.Ensembl38Id != "" || entity.Ensembl37Id != "" {
			updates = append(updates, entity)
		} else if unknownIds[entity.Id] != "" {
			entity.Ensembl38Id = unknownIds[entity.Id]
			updates = append(updates, entity)
		} else {
			log.Printf("Did not find Ensembl entry for %s", entity.Id)
		}
	}
	return
}

func (s Session) getEnsemblCoordUrl(gene string) string {
	return fmt.Sprintf("%s/lookup/id/%s?content-type=application/json;expand=1", s.getBuildUrl(), gene)
}

func (s Session) getCoordinates(id string) (coordinates EnsemblGeneObj, err error) {
	body, err := sendHttpRequest(s.getEnsemblCoordUrl(id))
	if err != nil {
		return
	}
	json.Unmarshal(body, &coordinates)
	return
}
