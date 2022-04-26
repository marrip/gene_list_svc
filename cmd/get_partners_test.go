package cmd

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/jarcoal/httpmock"
)

func TestGetPartners(t *testing.T) {
	var cases = map[string]struct {
		dbTableRow DbTableRow
		result     []DbTableRow
		wantErr    bool
	}{
		"Successful request": {
			DbTableRow{
				Id:    "ABL1",
				Class: "gene",
			},
			[]DbTableRow{
				{
					Id:    "GENE2",
					Class: "gene",
					Analyses: map[string]struct{}{
						"sv": struct{}{},
					},
				},
				{
					Id:    "GENE3",
					Class: "gene",
					Analyses: map[string]struct{}{
						"sv": struct{}{},
					},
				},
				{
					Id:    "ABL1",
					Class: "gene",
				},
			},
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "/gene-fusions/?id=1",
				httpmock.NewStringResponder(200, `ABL1 (1p1.1) GENE2 (1q1.1)</li><li class="border list-group-item">GENE3 (10p14) ABL1</ul>`))
			result, err := c.dbTableRow.getPartners()
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetIdUrl(t *testing.T) {
	var cases = map[string]struct {
		id      string
		result  string
		wantErr bool
	}{
		"Id exists": {
			"ABL1",
			"/gene-fusions/?id=1",
			false,
		},
		"Id is missing": {
			"GENE1",
			"",
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := getIdUrl(c.id)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetHtmlPage(t *testing.T) {
	var cases = map[string]struct {
		url     string
		result  []string
		wantErr bool
	}{
		"Successful request": {
			"http://atlasgeneticsoncology.org/gene-fusions/?id=1",
			[]string{"ABL1", "BCR"},
			false,
		},
		"Internal server error": {
			"http://atlasgeneticsoncology.org/gene-fusions/?id=2",
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "http://atlasgeneticsoncology.org/gene-fusions/?id=1",
				httpmock.NewStringResponder(200, "ABL1\nBCR"))
			httpmock.RegisterResponder("GET", "http://atlasgeneticsoncology.org/gene-fusions/?id=2",
				httpmock.NewStringResponder(500, ""))
			result, err := getHtmlPage(c.url)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestSendHttpRequest(t *testing.T) {
	var cases = map[string]struct {
		url     string
		result  []byte
		wantErr bool
	}{
		"Successful request": {
			"http://atlasgeneticsoncology.org/gene-fusions/?id=1",
			[]byte("<ul>"),
			false,
		},
		"Internal server error": {
			"http://atlasgeneticsoncology.org/gene-fusions/?id=2",
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", "http://atlasgeneticsoncology.org/gene-fusions/?id=1",
				httpmock.NewStringResponder(200, `<ul>`))
			httpmock.RegisterResponder("GET", "http://atlasgeneticsoncology.org/gene-fusions/?id=2",
				httpmock.NewStringResponder(500, ""))
			result, err := sendHttpRequest(c.url)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestScrapePartners(t *testing.T) {
	var cases = map[string]struct {
		lines  []string
		result []string
	}{
		"Scrape partners successfully": {
			[]string{`GENE1 (1p1.1) GENE2 (1q1.1)</li><li class="border list-group-item">GENE3 (10p14) GENE1</ul>`},
			[]string{"GENE1", "GENE2", "GENE3"},
		},
		"No partners scraped": {
			[]string{"<ul>"},
			nil,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := scrapePartners(c.lines)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGetPairs(t *testing.T) {
	var cases = map[string]struct {
		line   string
		result []string
	}{
		"Identify pairs successfully": {
			`GENE1 (1p1.1) GENE2 (1q1.1)</li><li class="border list-group-item">GENE3 (10p14) GENE1</ul>`,
			[]string{"GENE1", "GENE2", "GENE3", "GENE1"},
		},
		"No pairs identified": {
			" </ul>",
			nil,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := getPairs(c.line)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestExtractPartners(t *testing.T) {
	var cases = map[string]struct {
		line   string
		result []string
	}{
		"Identify pairs successfully": {
			"GENE1 (1p1.1) GENE2 (1q1.1)",
			[]string{"GENE1", "GENE2"},
		},
		"No pairs identified": {
			"</ul>",
			[]string{""},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := extractPartners(c.line)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	var cases = map[string]struct {
		slice  []string
		result []string
	}{
		"All strings are unique": {
			[]string{"GENE1", "GENE2", "GENE3"},
			[]string{"GENE1", "GENE2", "GENE3"},
		},
		"Some strings are duplicated": {
			[]string{"GENE1", "GENE2", "GENE1"},
			[]string{"GENE1", "GENE2"},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := unique(c.slice)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
