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
		if s.Build == "38" {
			gene, err = s.getCoordinates(entity.Ensembl38Id)
			if err != nil {
				log.Printf("Did not find coordinates for %s, %s", entity.Id, entity.Ensembl38Id)
				continue
			}
		} else if s.Build == "37" {
			gene, err = s.getCoordinates(entity.Ensembl37Id)
			if err != nil {
				log.Printf("Did not find coordinates for %s, %s", entity.Id, entity.Ensembl37Id)
				continue
			}
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

func sortRegions(regions []Region) (sortedRegions []Region) {
	sort.SliceStable(regions, func(i, j int) bool { return regions[i].End < regions[j].End })
	sort.SliceStable(regions, func(i, j int) bool { return regions[i].Start < regions[j].Start })
	sortedRegions = sortChromosomes(regions)
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

func getBedLine(regions []Region, prefix string) (line []string) {
	line = []string{
		fmt.Sprintf("%s%s", prefix, regions[0].Chromosome),
		strconv.Itoa(regions[0].Start),
		strconv.Itoa(regions[len(regions)-1].End),
		overlapAnnotation(regions),
	}
	return
}

func sortAndMergeRegions(regions []Region, prefix string) (lines [][]string, err error) {
	regions = sortRegions(regions)
	var overlapRegions []Region
	for i, region := range regions {
		if i == 0 {
			overlapRegions = []Region{region}
		} else if regionOverlap(regions[i-1], region) {
			overlapRegions = append(overlapRegions, region)
		} else {
			lines = append(lines, getBedLine(overlapRegions, prefix))
			overlapRegions = []Region{region}
		}
	}
	lines = append(lines, getBedLine(overlapRegions, prefix))
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
	var prefix string
	if s.Chr {
		prefix = "chr"
	}
	var lines [][]string
	lines, err = sortAndMergeRegions(regions, prefix)
	if err != nil {
		return
	}
	err = writeTsv(s.Bed, lines)
	return
}
