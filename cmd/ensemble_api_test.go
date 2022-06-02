package cmd

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/jarcoal/httpmock"
)

func TestGetEnsemblIds(t *testing.T) {
	var cases = map[string]struct {
		d       DbTableRow
		result  DbTableRow
		wantErr bool
	}{
		"Nothing to be done for region": {
			DbTableRow{
				Class: "region",
				Id:    "my_region",
			},
			DbTableRow{
				Class: "region",
				Id:    "my_region",
			},
			false,
		},
		"Set exon id for both genome build": {
			DbTableRow{
				Class: "exon",
				Id:    "ENSE0001",
			},
			DbTableRow{
				Class:       "exon",
				EnsemblId37: "ENSE0001",
				EnsemblId38: "ENSE0001",
				Id:          "ENSE0001",
			},
			false,
		},
		"Retrieve gene id": {
			DbTableRow{
				Class: "gene",
				Id:    "GENE1",
			},
			DbTableRow{
				Class:       "gene",
				EnsemblId37: "ENSG0001",
				EnsemblId38: "ENSG0001",
				Id:          "GENE1",
			},
			false,
		},
		"Internal server error": {
			DbTableRow{
				Class: "gene",
				Id:    "GENE2",
			},
			DbTableRow{
				Class: "gene",
				Id:    "GENE2",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "/xrefs/symbol/homo_sapiens/GENE1?content-type=application/json",
				httpmock.NewStringResponder(200, `[{"id": "ENSG0001"}]`))
			httpmock.RegisterResponder("GET", "/lookup/id/ENSG0001?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "1"}`))
			httpmock.RegisterResponder("GET", "/lookup/id/GENE2?content-type=application/json",
				httpmock.NewStringResponder(500, ""))
			err := c.d.getEnsemblIds()
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(c.d, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetEnsemblId(t *testing.T) {
	var cases = map[string]struct {
		d       DbTableRow
		build   string
		result  string
		wantErr bool
	}{
		"Retrieve gene id successfully": {
			DbTableRow{
				Class: "gene",
				Id:    "GENE1",
			},
			"38",
			"ENSG0001",
			false,
		},
		"Retrieve transcript id successfully": {
			DbTableRow{
				Class: "transcript",
				Id:    "TRANSCRIPT1",
			},
			"38",
			"ENST0001",
			false,
		},
		"Internal server error": {
			DbTableRow{
				Class: "gene",
				Id:    "GENE2",
			},
			"38",
			"",
			true,
		},
		"Chromosome is invalid": {
			DbTableRow{
				Class: "gene",
				Id:    "GENE3",
			},
			"38",
			"",
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "/xrefs/symbol/homo_sapiens/GENE1?content-type=application/json",
				httpmock.NewStringResponder(200, `[{"id": "ENSG0001"}]`))
			httpmock.RegisterResponder("GET", "/lookup/id/ENSG0001?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "1"}`))
			httpmock.RegisterResponder("GET", "/xrefs/symbol/homo_sapiens/TRANSCRIPT1?content-type=application/json",
				httpmock.NewStringResponder(200, `[{"id": "ENST0001"}]`))
			httpmock.RegisterResponder("GET", "/lookup/id/ENST0001?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "1"}`))
			httpmock.RegisterResponder("GET", "/lookup/id/GENE2?content-type=application/json",
				httpmock.NewStringResponder(500, ""))
			httpmock.RegisterResponder("GET", "/lookup/id/GENE3?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "100"}`))
			result, err := c.d.getEnsemblId(c.build)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetCrossRefUrl(t *testing.T) {
	var cases = map[string]struct {
		d      DbTableRow
		build  string
		result string
	}{
		"Retrieve url sucessfully": {
			DbTableRow{
				Id: "GENE1",
			},
			"38",
			"/xrefs/symbol/homo_sapiens/GENE1?content-type=application/json",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := c.d.getCrossRefUrl(c.build)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetBuildUrl(t *testing.T) {
	var cases = map[string]struct {
		build  string
		result string
	}{
		"Genome build is 38": {
			"38",
			"my.ensembl38",
		},
		"Genome build is 37": {
			"37",
			"my.ensembl37",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			session = Session{
				Web: web{
					Ensembl37: "my.ensembl37",
					Ensembl38: "my.ensembl38",
				},
			}
			result := getBuildUrl(c.build)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
			session = Session{}
		})
	}
}

func TestCheckEnsemblIdChromosome(t *testing.T) {
	var cases = map[string]struct {
		id      string
		build   string
		result  bool
		wantErr bool
	}{
		"Chromosome is valid": {
			"GENE1",
			"38",
			true,
			false,
		},
		"Internal server error": {
			"GENE2",
			"38",
			false,
			true,
		},
		"Chromosome is not valid": {
			"GENE3",
			"38",
			false,
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "/lookup/id/GENE1?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "1"}`))
			httpmock.RegisterResponder("GET", "/lookup/id/GENE2?content-type=application/json",
				httpmock.NewStringResponder(500, ""))
			httpmock.RegisterResponder("GET", "/lookup/id/GENE3?content-type=application/json",
				httpmock.NewStringResponder(200, `{"seq_region_name": "100"}`))
			result, err := checkEnsemblIdChromosome(c.id, c.build)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetLookUpUrl(t *testing.T) {
	var cases = map[string]struct {
		id     string
		build  string
		expand bool
		result string
	}{
		"Url contains expand": {
			"GENE1",
			"38",
			true,
			"/lookup/id/GENE1?content-type=application/json;expand=1",
		},
		"Url without expand": {
			"GENE1",
			"38",
			false,
			"/lookup/id/GENE1?content-type=application/json",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := getLookUpUrl(c.id, c.build, c.expand)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
