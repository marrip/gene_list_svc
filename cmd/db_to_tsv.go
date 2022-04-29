package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func dbToTsv() (err error) {
	rows, err := session.Db.Connection.getRegions()
	if err != nil {
		return
	}
	var regions []DbTableRow
	for _, row := range rows {
		if row.Class == "region" {
			regions = append(regions, row)
			continue
		}
		switch session.Analysis {
		case "pindel", "sv":
			err = row.getCompleteRegion()
			if err != nil {
				continue
			}
			regions = append(regions, row)
		case "cnv", "snv":
			if row.Class == "exon" {
				err = row.getCompleteRegion()
				if err != nil {
					continue
				}
				regions = append(regions, row)
			} else {
			}
		}
	}
	fmt.Println(regions)
	fmt.Println(len(regions))
	return
}

func (d *DbTableRow) getCompleteRegion() (err error) {
	body, err := d.getCoordinates(false)
	if err != nil {
		return
	}
	var obj EnsemblBaseObj
	json.Unmarshal(body, &obj)
	d.Chromosome = obj.Chromosome
	d.Start = strconv.Itoa(obj.Start)
	d.End = strconv.Itoa(obj.End)
	return
}

func (d DbTableRow) getExons() (regions []DbTableRow, err error) {
	body, err := d.getCoordinates(true)
	if err != nil {
		return
	}
	var exons []EnsemblBaseObj
	if d.Class == "gene" {
		var obj EnsemblGeneObj
		json.Unmarshal(body, &obj)
		for _, transcript := range obj.Transcripts {
			for _, exon := range transcript.Exons {
				exon.Transcript = []string{transcript.EnsemblId}
				exons = append(exons, exon)
			}
		}
	} else if d.Class == "transcript" {
		var obj EnsemblTransObj
		json.Unmarshal(body, &obj)
		exons = obj.Exons
	}
	exons = uniqueExons(exons)
	return
}

func uniqueExons(exons []EnsemblBaseObj) (uniqueExons []EnsemblBaseObj) {
	keys := make(map[string]bool)
	for _, exon := range exons {
		if !keys[exon.EnsemblId] {
			keys[exon.EnsemblId] = true
			uniqueExons = append(uniqueExons, exon)
		}
	}
	return
}
