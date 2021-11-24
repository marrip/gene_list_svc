package main

var analyses = []string{
	"snv",
	"cnv",
	"sv",
	"pindel",
}

var lists = []string{
	"aml",
	"aml_ext",
	"all",
}

var chromosomes = map[string]bool{
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
}

var unknownIds = map[string]string{
	"AFND":           "ENSG00000130396",
	"RP11-1407O15.2": "ENSG00000174093",
	"SEPT5-GP1BB":    "ENSG00000284874",
	"TMX2-CTNND1":    "ENSG00000288534",
}
