package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func (d *DbTableRow) getEnsemblIds() (err error) {
	if d.Class == "region" {
		return
	} else if d.Class == "exon" {
		d.EnsemblId38 = d.Id
		d.EnsemblId37 = d.Id
	} else {
		if d.EnsemblId38, err = d.getEnsemblId("38"); err != nil {
			return
		}
		if d.EnsemblId37, err = d.getEnsemblId("37"); err != nil {
			return
		}
	}
	return
}

func (d *DbTableRow) getEnsemblId(build string) (id string, err error) {
	geneRegex := regexp.MustCompile("ENSG")
	transRegex := regexp.MustCompile("ENST")
	body, err := sendHttpRequest(d.getCrossRefUrl(build))
	if err != nil {
		return
	}
	var jsonObj []EnsemblGeneObj
	json.Unmarshal(body, &jsonObj)
	for _, element := range jsonObj {
		if d.Class == "gene" && geneRegex.MatchString(element.EnsemblId) {
			var valid bool
			valid, err = checkEnsemblIdChromosome(element.EnsemblId, build)
			if err != nil {
				return
			}
			if valid {
				id = element.EnsemblId
				return
			}
		} else if d.Class == "transcript" && transRegex.MatchString(element.EnsemblId) {
			id = element.EnsemblId
			return
		}
	}
	return
}

func (d DbTableRow) getCrossRefUrl(build string) string {
	return fmt.Sprintf("%s/xrefs/symbol/homo_sapiens/%s?content-type=application/json", getBuildUrl(build), d.Id)
}

func getBuildUrl(build string) string {
	if build == "37" {
		return session.Web.Ensembl37
	}
	return session.Web.Ensembl38
}

func checkEnsemblIdChromosome(id string, build string) (valid bool, err error) {
	body, err := sendHttpRequest(getLookUpUrl(id, build, false))
	if err != nil {
		return
	}
	var jsonObj EnsemblGeneObj
	json.Unmarshal(body, &jsonObj)
	chromosomes := generateChromosomeMap()
	valid = chromosomes[jsonObj.Chromosome]
	return
}

func getLookUpUrl(id string, build string, expand bool) string {
	var expandString string
	if expand {
		expandString = ";expand=1"
	}
	return fmt.Sprintf("%s/lookup/id/%s?content-type=application/json%s", getBuildUrl(build), id, expandString)
}

//func (s Session) getCoordinates(id string) (coordinates EnsemblGeneObj, err error) {
//	body, err := sendHttpRequest(s.getEnsemblCoordUrl(id))
//	if err != nil {
//		return
//	}
//	json.Unmarshal(body, &coordinates)
//	return
//}
