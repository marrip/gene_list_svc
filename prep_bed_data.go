package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

func (s Session) collectExonCoordinates(entities []Entity) (regions []Region, err error) {
	for _, entity := range entities {
		var gene EnsemblGeneObj
		log.Printf("Retrieve coordinates for %s", entity.Id)
		gene, err = s.getCoordinates(entity.Ensembl38Id)
		if err != nil {
			log.Printf("Did not find coordinates for %s", entity.Ensembl38Id)
			continue
		}
		for _, transcript := range gene.Transcripts {
			for _, exon := range transcript.Exons {
				regions = append(regions, Region{
					Gene:       gene.Id,
					Id:         exon.EnsemblId,
					Chromosome: exon.Chromosome,
					Start:      exon.Start - 10,
					End:        exon.End + 10,
				})
			}
		}
	}
	return
}

func uniqueExons(exons []Region) (uniqueExons []Region) {
	keys := make(map[string]bool)
	for _, exon := range exons {
		if !keys[exon.Id] {
			keys[exon.Id] = true
			uniqueExons = append(uniqueExons, exon)
		}
	}
	return
}

func (s Session) collectGeneCoordinates(entities []Entity) (regions []Region, err error) {
	for _, entity := range entities {
		var gene EnsemblGeneObj
		log.Printf("Retrieve coordinates for %s", entity.Id)
		gene, err = s.getCoordinates(entity.Ensembl38Id)
		if err != nil {
			log.Printf("Did not find coordinates for %s", entity.Ensembl38Id)
			continue
		}
		regions = append(regions, Region{
			Gene:       gene.Id,
			Id:         gene.EnsemblId,
			Chromosome: gene.Chromosome,
			Start:      gene.Start - 50,
			End:        gene.End + 50,
		})
	}
	return
}

func sortChromosomes(regions []Region) (sortedRegions []Region) {
	for _, chromosome := range chromosomes {
		for _, region := range regions {
			if region.Chromosome == chromosome {
				sortedRegions = append(sortedRegions, region)
			}
		}
	}
	return
}

func regionOverlap(a, b Region) bool {
	if a.Chromosome == b.Chromosome && a.End >= b.Start {
		return true
	}
	return false
}

func overlapAnnotation(regions []Region) (annotation string) {
	for i, region := range regions {
		if i == 0 {
			annotation = fmt.Sprintf("%s|%s", region.Gene, region.Id)
		} else if regions[i-1].Gene == region.Gene {
			annotation = fmt.Sprintf("%s|%s", annotation, region.Id)
		} else {
			annotation = fmt.Sprintf("%s,%s|%s", annotation, region.Gene, region.Id)
		}
	}
	return
}

func sortAndMergeRegions(regions []Region) (lines [][]string, err error) {
	sort.SliceStable(regions, func(i, j int) bool { return regions[i].Start < regions[j].Start })
	regions = sortChromosomes(regions)
	var overlapRegions []Region
	for i, region := range regions {
		if i == 0 {
			overlapRegions = []Region{region}
		} else if regionOverlap(regions[i-1], region) {
			overlapRegions = append(overlapRegions, region)
		} else {
			lines = append(lines, []string{overlapRegions[0].Chromosome, strconv.Itoa(overlapRegions[0].Start), strconv.Itoa(overlapRegions[len(overlapRegions)-1].End), overlapAnnotation(overlapRegions)})
			overlapRegions = []Region{region}
		}
	}
	return
}

func (s Session) getCoordinatesforEntities() (err error) {
	var entities []Entity
	entities, err = s.getEntityList()
	if err != nil {
		return
	}
	var regions []Region
	switch s.Analysis {
	case "snv":
		regions, err = s.collectExonCoordinates(entities)
		if err != nil {
			return
		}
		regions = uniqueExons(regions)
	case "cnv":
		log.Println("Running cnv")
	case "sv":
		regions, err = s.collectGeneCoordinates(entities)
		if err != nil {
			return
		}
	case "pindel":
		regions, err = s.collectGeneCoordinates(entities)
		if err != nil {
			return
		}
	}
	var lines [][]string
	lines, err = sortAndMergeRegions(regions)
	if err != nil {
		return
	}
	err = writeTsv(s.Bed, lines)
	return
}
