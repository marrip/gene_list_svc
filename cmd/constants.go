package cmd

var analyses = map[string]struct{}{
	"snv":    {},
	"cnv":    {},
	"sv":     {},
	"pindel": {},
}

var classes = map[string]struct{}{
	"gene":       {},
	"transcript": {},
	"exon":       {},
	"region":     {},
}

var tsvHeader = map[string]bool{
	"analyses":         true,
	"class":            true,
	"coordinates":      false,
	"id":               true,
	"include_partners": true,
	"tables":           true,
}
