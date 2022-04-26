package cmd

var analyses = map[string]struct{}{
	"snv":    {},
	"cnv":    {},
	"sv":     {},
	"pindel": {},
}

var atlasIds = map[string]string{
	"ABL1":   "1",
	"ABL2":   "226",
	"CRLF2":  "51262",
	"CSF1R":  "40161",
	"ETV6":   "38",
	"FGFR1":  "113",
	"IGH":    "40",
	"JAK2":   "98",
	"KAT6A":  "25",
	"KMT2A":  "13",
	"MLLT10": "4",
	"NUP98":  "63",
	"NUTM1":  "41595",
	"PDGFRB": "21",
	"RARA":   "46",
	"RUNX1":  "52",
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
