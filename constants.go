package main

var analyses = []string{
	"snv",
	"cnv",
	"sv",
	"pindel",
}

var builds = []string{
	"37",
	"38",
}

var lists = []string{
	"aml",
	"aml_ext",
	"all",
}

var chromosomes = []string{
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"10",
	"11",
	"12",
	"13",
	"14",
	"15",
	"16",
	"17",
	"18",
	"19",
	"20",
	"21",
	"22",
	"X",
	"Y",
}

var unknownIds = map[string]string{
	"AFND":        "ENSG00000130396",
	"CSFR3":       "ENSG00000119535",
	"SEPT5-GP1BB": "ENSG00000284874",
	"TMX2-CTNND1": "ENSG00000288534",
}
