package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func dbToTsv() (err error) {
	rows, err := session.Db.Connection.getRegions()
	if err != nil {
		return
	}
	var regions []EnsemblBaseObj
	switch session.Analysis {
	case "pindel", "sv":
		regions, err = prepForPindelSv(rows)
		if err != nil {
			return
		}
	case "cnv", "snv":
		regions, err = prepForCnvSnv(rows)
		if err != nil {
			return
		}
	}
	if err = regionsToTsv(regions); err != nil {
		return
	}
	return
}

func prepForPindelSv(rows []DbTableRow) (regions []EnsemblBaseObj, err error) {
	for _, row := range rows {
		var region EnsemblBaseObj
		if row.Class == "region" {
			region, err = row.rowToRegion()
			if err != nil {
				return
			}
		} else {
			region, err = row.getCompleteRegion(50)
			if err != nil {
				return
			}
		}
		regions = append(regions, region)
	}
	return
}

func (d DbTableRow) rowToRegion() (region EnsemblBaseObj, err error) {
	region.Annotation = d.Id
	region.Chromosome = d.Chromosome
	if region.Start, err = strconv.Atoi(d.Start); err != nil {
		return
	}
	if region.End, err = strconv.Atoi(d.End); err != nil {
		return
	}
	return
}

func (d DbTableRow) getCompleteRegion(size int) (region EnsemblBaseObj, err error) {
	body, err := d.getCoordinates(false)
	if err != nil {
		return
	}
	json.Unmarshal(body, &region)
	if d.Id != region.EnsemblId {
		region.Annotation = fmt.Sprintf("%s|%s", d.Id, region.EnsemblId)
	} else {
		region.Annotation = d.Id
	}
	region.addWindow(size)
	return
}

func (o *EnsemblBaseObj) addWindow(size int) {
	o.Start = o.Start - size
	o.End = o.End + size
}

func prepForCnvSnv(rows []DbTableRow) (regions []EnsemblBaseObj, err error) {
	for _, row := range rows {
		var region EnsemblBaseObj
		if row.Class == "region" {
			region, err = row.rowToRegion()
			if err != nil {
				return
			}
			regions = append(regions, region)
		} else if row.Class == "exon" {
			region, err = row.getCompleteRegion(10)
			if err != nil {
				return
			}
			regions = append(regions, region)
		} else {
			var exons []EnsemblBaseObj
			exons, err = row.getExons(10)
			if err != nil {
				return
			}
			regions = append(regions, exons...)
		}
	}
	return
}

func (d DbTableRow) getExons(size int) (regions []EnsemblBaseObj, err error) {
	body, err := d.getCoordinates(true)
	if err != nil {
		return
	}
	if d.Class == "gene" {
		var obj EnsemblGeneObj
		json.Unmarshal(body, &obj)
		var exons []EnsemblBaseObj
		for _, transcript := range obj.Transcripts {
			for _, exon := range transcript.Exons {
				exon.Transcript = transcript.EnsemblId
				exons = append(exons, exon)
			}
		}
		geneAnnotation := fmt.Sprintf("%s|%s", obj.Id, obj.EnsemblId)
		regions = append(regions, uniqueExons(geneAnnotation, exons)...)
	} else if d.Class == "transcript" {
		var obj EnsemblTransObj
		json.Unmarshal(body, &obj)
		for _, exon := range obj.Exons {
			exon.addWindow(size)
			exon.Annotation = fmt.Sprintf("%s|%s|%s", d.Id, obj.EnsemblId, exon.EnsemblId)
			regions = append(regions, exon)
		}
	}
	return
}

func uniqueExons(annotation string, exons []EnsemblBaseObj) (uniqueExons []EnsemblBaseObj) {
	uniqs := make(map[string]EnsemblBaseObj)
	for _, exon := range exons {
		if _, present := uniqs[exon.EnsemblId]; !present {
			exon.addWindow(10)
			exon.Annotation = fmt.Sprintf("%s|%s|%s", annotation, exon.Transcript, exon.EnsemblId)
			uniqs[exon.EnsemblId] = exon
		} else {
			uniq := uniqs[exon.EnsemblId]
			exonAnnotation := strings.Split(uniq.Annotation, "|")
			uniq.Annotation = fmt.Sprintf("%s|%s|%s&%s|%s", exonAnnotation[0], exonAnnotation[1], exonAnnotation[2], exon.Transcript, exonAnnotation[3])
			uniqs[exon.EnsemblId] = uniq
		}
	}
	for _, uniq := range uniqs {
		uniqueExons = append(uniqueExons, uniq)
	}
	return
}

func regionsToTsv(regions []EnsemblBaseObj) (err error) {
	regions = sortRegions(regions)
	lines := regionsToSlices(regions)
	err = writeTsv(session.Bed, lines)
	return
}

func sortRegions(regions []EnsemblBaseObj) (sortedRegions []EnsemblBaseObj) {
	sort.SliceStable(regions, func(i, j int) bool { return regions[i].End < regions[j].End })
	sort.SliceStable(regions, func(i, j int) bool { return regions[i].Start < regions[j].Start })
	sortedRegions = sortChromosomes(regions)
	return
}

func sortChromosomes(regions []EnsemblBaseObj) (sortedRegions []EnsemblBaseObj) {
	chromosomes := generateChromosomeSlice()
	for _, chromosome := range chromosomes {
		for _, region := range regions {
			if region.Chromosome == chromosome {
				sortedRegions = append(sortedRegions, region)
			}
		}
	}
	return
}

func generateChromosomeSlice() (chromosomes []string) {
	for i := 1; i <= 22; i++ {
		chromosomes = append(chromosomes, strconv.Itoa(i))
	}
	for _, chromosome := range []string{"X", "Y", "M"} {
		chromosomes = append(chromosomes, chromosome)
	}
	return
}

func regionsToSlices(regions []EnsemblBaseObj) (lines [][]string) {
	var overlapRegions []EnsemblBaseObj
	for i, region := range regions {
		if i == 0 {
			overlapRegions = []EnsemblBaseObj{region}
		} else if identifyOverlap(regions[i-1], region) {
			overlapRegions = append(overlapRegions, region)
		} else {
			lines = append(lines, mergeOverlappingRegions(overlapRegions))
			overlapRegions = []EnsemblBaseObj{region}
		}
	}
	return
}

func identifyOverlap(a, b EnsemblBaseObj) bool {
	if a.Chromosome == b.Chromosome && a.End >= b.Start {
		return true
	}
	return false
}

func mergeOverlappingRegions(regions []EnsemblBaseObj) (line []string) {
	var chromosome, annotation string
	var start, end int
	for i, region := range regions {
		if i == 0 {
			start = region.Start
			end = region.End
			annotation = region.Annotation
		} else {
			annotation = fmt.Sprintf("%s;%s", annotation, region.Annotation)
		}
		if region.Start < start {
			start = region.Start
		}
		if region.End > end {
			end = region.End
		}
	}
	if session.Chr {
		chromosome = fmt.Sprintf("chr%s", regions[0].Chromosome)
	} else {
		chromosome = regions[0].Chromosome
	}
	line = []string{chromosome, strconv.Itoa(start), strconv.Itoa(end), annotation}
	return
}
