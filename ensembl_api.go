package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

func (s Session) getEnsemblIdUrl(gene string) string {
	return fmt.Sprintf("%s/xrefs/symbol/homo_sapiens/%s?content-type=application/json", s.EnsemblRestUrl, gene)
}

func (s Session) getEnsemblSmallGeneUrl(id string) string {
	return fmt.Sprintf("%s/lookup/id/%s?content-type=application/json", s.EnsemblRestUrl, id)
}

func (s Session) checkEnsemblIdChromosome(id string) (isMain bool, err error) {
	body, err := sendHttpRequest(s.getEnsemblSmallGeneUrl(id))
	if err != nil {
		return
	}
	var jsonObj EnsemblGeneObj
	json.Unmarshal(body, &jsonObj)
	isMain = chromosomes[jsonObj.Chromosome]
	return
}

func (s Session) getEnsemblId(gene string) (id string, err error) {
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
			isMain, err = s.checkEnsemblIdChromosome(element.EnsemblId)
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
		entity.EnsemblId, err = s.getEnsemblId(entity.Id)
		if err != nil {
			return
		}
		if entity.EnsemblId != "" {
			updates = append(updates, entity)
		} else if unknownIds[entity.Id] != "" {
			entity.EnsemblId = unknownIds[entity.Id]
			updates = append(updates, entity)
		} else {
			log.Printf("Did not find Ensembl entry for %s", entity.Id)
		}
	}
	return
}

func (s Session) getEnsemblCoordUrl(gene string) string {
	return fmt.Sprintf("%s/lookup/id/%s?content-type=application/json;expand=1", s.EnsemblRestUrl, gene)
}

func (s Session) getCoordinates(id string) (coordinates EnsemblGeneObj, err error) {
	body, err := sendHttpRequest(s.getEnsemblCoordUrl(id))
	if err != nil {
		return
	}
	json.Unmarshal(body, &coordinates)
	return
}
