package main

import (
	"testing"

	"github.com/go-test/deep"
)

func TestUniqueExons(t *testing.T) {
	var cases = map[string]struct {
		exons  []Region
		result []Region
	}{
		"All exons unique": {
			[]Region{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
				{
					Id: "3",
				},
			},
			[]Region{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
				{
					Id: "3",
				},
			},
		},
		"All exons identical": {
			[]Region{
				{
					Id: "1",
				},
				{
					Id: "1",
				},
				{
					Id: "1",
				},
			},
			[]Region{
				{
					Id: "1",
				},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := uniqueExons(c.exons)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestSortChromosomes(t *testing.T) {
	var cases = map[string]struct {
		regions []Region
		result  []Region
	}{
		"Successfully sort all chromosomes": {
			[]Region{
				{
					Chromosome: "2",
				},
				{
					Chromosome: "1",
				},
				{
					Chromosome: "2",
				},
			},
			[]Region{
				{
					Chromosome: "1",
				},
				{
					Chromosome: "2",
				},
				{
					Chromosome: "2",
				},
			},
		},
		"Remove non-existent chromosomes": {
			[]Region{
				{
					Chromosome: "Z",
				},
				{
					Chromosome: "1",
				},
			},
			[]Region{
				{
					Chromosome: "1",
				},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := sortChromosomes(c.regions)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestSortRegions(t *testing.T) {
	var cases = map[string]struct {
		regions []Region
		result  []Region
	}{
		"Successfully sort regions": {
			[]Region{
				{
					Chromosome: "1",
					Start:      200,
					End:        220,
				},
				{
					Chromosome: "1",
					Start:      100,
					End:        230,
				},
			},
			[]Region{
				{
					Chromosome: "1",
					Start:      100,
					End:        230,
				},
				{
					Chromosome: "1",
					Start:      200,
					End:        220,
				},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := sortRegions(c.regions)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestRegionOverlap(t *testing.T) {
	var cases = map[string]struct {
		a      Region
		b      Region
		result bool
	}{
		"Regions overlap": {
			Region{
				Chromosome: "1",
				End:        120,
			},
			Region{
				Chromosome: "1",
				Start:      100,
			},
			true,
		},
		"Regions on different chromosomes": {
			Region{
				Chromosome: "1",
				End:        120,
			},
			Region{
				Chromosome: "2",
				Start:      100,
			},
			false,
		},
		"Regions on same chromosomes but do not overlap": {
			Region{
				Chromosome: "1",
				End:        10,
			},
			Region{
				Chromosome: "2",
				Start:      100,
			},
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := regionOverlap(c.a, c.b)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestOverlapAnnotation(t *testing.T) {
	var cases = map[string]struct {
		regions []Region
		result  string
	}{
		"Single region annotated": {
			[]Region{
				{
					Id:   "1",
					Gene: "Gene1",
				},
			},
			"Gene1|1",
		},
		"Combined annotation for same gene": {
			[]Region{
				{
					Id:   "1",
					Gene: "Gene1",
				},
				{
					Id:   "2",
					Gene: "Gene1",
				},
			},
			"Gene1|1|2",
		},
		"Combined annotation for different genes": {
			[]Region{
				{
					Id:   "1",
					Gene: "Gene1",
				},
				{
					Id:   "1",
					Gene: "Gene2",
				},
			},
			"Gene1|1,Gene2|1",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := overlapAnnotation(c.regions)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetBedLine(t *testing.T) {
	var cases = map[string]struct {
		regions []Region
		result  []string
	}{
		"Single region bed line": {
			[]Region{
				{
					Id:         "1",
					Gene:       "Gene1",
					Chromosome: "1",
					Start:      0,
					End:        10,
				},
			},
			[]string{
				"chr1",
				"0",
				"10",
				"Gene1|1",
			},
		},
		"Bed line for overlapping regions": {
			[]Region{
				{
					Id:         "1",
					Gene:       "Gene1",
					Chromosome: "1",
					Start:      0,
					End:        10,
				},
				{
					Id:         "2",
					Gene:       "Gene1",
					Chromosome: "1",
					Start:      8,
					End:        25,
				},
			},
			[]string{
				"chr1",
				"0",
				"25",
				"Gene1|1|2",
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := getBedLine(c.regions, "chr")
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestSortAndMergeRegions(t *testing.T) {
	var cases = map[string]struct {
		regions []Region
		result  [][]string
		wantErr bool
	}{
		"Succesfully sort and annotate": {
			[]Region{
				{
					Id:         "1",
					Gene:       "Gene4",
					Chromosome: "2",
					Start:      1,
					End:        15,
				},
				{
					Id:         "1",
					Gene:       "Gene2",
					Chromosome: "1",
					Start:      9,
					End:        15,
				},
				{
					Id:         "1",
					Gene:       "Gene3",
					Chromosome: "1",
					Start:      20,
					End:        30,
				},
				{
					Id:         "1",
					Gene:       "Gene1",
					Chromosome: "1",
					Start:      1,
					End:        10,
				},
			},
			[][]string{
				{
					"chr1",
					"1",
					"15",
					"Gene1|1,Gene2|1",
				},
				{
					"chr1",
					"20",
					"30",
					"Gene3|1",
				},
				{
					"chr2",
					"1",
					"15",
					"Gene4|1",
				},
			},
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := sortAndMergeRegions(c.regions, "chr")
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
