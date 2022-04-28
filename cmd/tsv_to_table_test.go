package cmd

import (
	"testing"

	"github.com/go-test/deep"
)

func TestValidateTsvHeader(t *testing.T) {
	var cases = map[string]struct {
		header  []string
		wantErr bool
	}{
		"Header is complete": {
			[]string{
				"analyses",
				"class",
				"coordinates",
				"id",
				"include_partners",
				"tables",
			},
			false,
		},
		"Header is missing coordinates": {
			[]string{
				"analyses",
				"class",
				"id",
				"include_partners",
				"tables",
			},
			false,
		},
		"Header is missing required columns": {
			[]string{
				"class",
				"id",
				"tables",
			},
			true,
		},
		"Header is complete but has extra columns": {
			[]string{
				"analyses",
				"class",
				"id",
				"include_partners",
				"nonsense",
				"tables",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateTsvHeader(c.header)
			checkError(t, err, c.wantErr)
		})
	}
}

func TestRowToMap(t *testing.T) {
	var cases = map[string]struct {
		row     []string
		header  []string
		result  map[string]string
		wantErr bool
	}{
		"Header and row have matching length": {
			[]string{
				"gene",
				"RUNX1",
			},
			[]string{
				"class",
				"id",
			},
			map[string]string{
				"class": "gene",
				"id":    "RUNX1",
			},
			false,
		},
		"Header and row length do not match": {
			[]string{
				"gene",
				"RUNX1",
			},
			[]string{
				"analyses",
				"class",
				"id",
			},
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := rowToMap(c.row, c.header)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestMapToDbRow(t *testing.T) {
	var cases = map[string]struct {
		row     map[string]string
		result  DbTableRow
		wantErr bool
	}{
		"All values are valid": {
			map[string]string{
				"analyses":         "snv,sv",
				"class":            "gene",
				"id":               "RUNX1",
				"include_partners": "true",
				"tables":           "my_list",
			},
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": struct{}{},
					"sv":  struct{}{},
				},
				Class:           "gene",
				Id:              "RUNX1",
				IncludePartners: true,
				Tables: []string{
					"my_list",
				},
			},
			false,
		},
		"Region is class region": {
			map[string]string{
				"analyses":         "snv",
				"class":            "region",
				"coordinates":      "chr1:1-10",
				"id":               "REGION1",
				"include_partners": "false",
				"tables":           "my_list",
			},
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": struct{}{},
				},
				Chromosome:      "1",
				Class:           "region",
				End:             "10",
				Id:              "REGION1",
				IncludePartners: false,
				Start:           "1",
				Tables: []string{
					"my_list",
				},
			},
			false,
		},
		"Analysis does not exist": {
			map[string]string{
				"analyses": "tsv",
			},
			DbTableRow{
				Analyses: map[string]struct{}{},
			},
			true,
		},
		"Class does not exist": {
			map[string]string{
				"analyses": "sv",
				"class":    "nonesense",
			},
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
			},
			true,
		},
		"Coordinates are invalid": {
			map[string]string{
				"analyses":    "sv",
				"class":       "region",
				"coordinates": "300:1-99",
			},
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Class: "region",
			},
			true,
		},
		"Analyses do not allow include partners": {
			map[string]string{
				"analyses":         "snv",
				"class":            "gene",
				"include_partners": "true",
			},
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": struct{}{},
				},
				Class: "gene",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := mapToDbRow(c.row)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestValidateAnalyses(t *testing.T) {
	var cases = map[string]struct {
		dbTableRow DbTableRow
		analyses   string
		result     DbTableRow
		wantErr    bool
	}{
		"Analyses exist": {
			DbTableRow{},
			"snv,sv",
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": {},
					"sv":  {},
				},
			},
			false,
		},
		"Analyses do not exist": {
			DbTableRow{},
			"snv,nonsense",
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": {},
				},
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.dbTableRow.validateAnalyses(c.analyses)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.dbTableRow, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestValidateClass(t *testing.T) {
	var cases = map[string]struct {
		dbTableRow DbTableRow
		class      string
		result     DbTableRow
		wantErr    bool
	}{
		"Class exists": {
			DbTableRow{},
			"gene",
			DbTableRow{
				Class: "gene",
			},
			false,
		},
		"Class does not exist": {
			DbTableRow{},
			"nonsense",
			DbTableRow{},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.dbTableRow.validateClass(c.class)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.dbTableRow, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	var cases = map[string]struct {
		dbTableRow  DbTableRow
		coordinates string
		result      DbTableRow
		wantErr     bool
	}{
		"Coordinates are valid without prefix": {
			DbTableRow{
				Class: "region",
			},
			"1:0-10",
			DbTableRow{
				Chromosome: "1",
				Class:      "region",
				End:        "10",
				Start:      "0",
			},
			false,
		},
		"Coordinates are valid with prefix": {
			DbTableRow{
				Class: "region",
			},
			"chrX:10-24",
			DbTableRow{
				Chromosome: "X",
				Class:      "region",
				End:        "24",
				Start:      "10",
			},
			false,
		},
		"Coordinates linked to wrong class": {
			DbTableRow{
				Class: "gene",
			},
			"chr10:10-24",
			DbTableRow{
				Class: "gene",
			},
			false,
		},
		"Coordinates are invalid": {
			DbTableRow{
				Class: "region",
			},
			"chr100:10-24",
			DbTableRow{
				Class: "region",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.dbTableRow.validateCoordinates(c.coordinates)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.dbTableRow, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestValidateChromosome(t *testing.T) {
	var cases = map[string]struct {
		dbTableRow DbTableRow
		chromosome string
		result     DbTableRow
		wantErr    bool
	}{
		"Chromosome is valid": {
			DbTableRow{},
			"1",
			DbTableRow{
				Chromosome: "1",
			},
			false,
		},
		"Chromosome is invalid": {
			DbTableRow{},
			"99",
			DbTableRow{},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.dbTableRow.validateChromosome(c.chromosome)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.dbTableRow, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGenerateChromosomeMap(t *testing.T) {
	var cases = map[string]struct {
		result map[string]bool
	}{
		"Sucessfully generate chromosome map": {
			map[string]bool{
				"1":  true,
				"2":  true,
				"3":  true,
				"4":  true,
				"5":  true,
				"6":  true,
				"7":  true,
				"8":  true,
				"9":  true,
				"10": true,
				"11": true,
				"12": true,
				"13": true,
				"14": true,
				"15": true,
				"16": true,
				"17": true,
				"18": true,
				"19": true,
				"20": true,
				"21": true,
				"22": true,
				"X":  true,
				"Y":  true,
				"M":  true,
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := generateChromosomeMap()
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestValidateIncludePartners(t *testing.T) {
	var cases = map[string]struct {
		row     DbTableRow
		include string
		result  DbTableRow
		wantErr bool
	}{
		"Sucessfully parse value": {
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Class: "gene",
			},
			"true",
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Class:           "gene",
				IncludePartners: true,
			},
			false,
		},
		"Value can not be parsed": {
			DbTableRow{},
			"nonesense",
			DbTableRow{},
			true,
		},
		"Analysis is not sv": {
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": struct{}{},
				},
				Class: "gene",
			},
			"true",
			DbTableRow{
				Analyses: map[string]struct{}{
					"snv": struct{}{},
				},
				Class: "gene",
			},
			true,
		},
		"Class is not gene": {
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Class: "region",
			},
			"true",
			DbTableRow{
				Analyses: map[string]struct{}{
					"sv": struct{}{},
				},
				Class: "region",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.row.validateIncludePartners(c.include)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.row, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
