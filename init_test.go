package main

import (
	"flag"
	"os"
	"testing"

	"github.com/go-test/deep"
)

func checkError(t *testing.T, err error, exp bool) {
	if (err != nil) != exp {
		t.Errorf("Expectation and result are different. Error is\n%v", err)
	}
}

func TestSliceContainsString(t *testing.T) {
	var cases = map[string]struct {
		arg    string
		slice  []string
		result bool
	}{
		"String is in slice": {
			"aml",
			[]string{"aml", "all"},
			true,
		},
		"String is not in slice": {
			"aml",
			[]string{"all"},
			false,
		},
		"Slice is empty": {
			"aml",
			[]string{},
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := sliceContainsString(c.arg, c.slice)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestCheckListAndAnalysis(t *testing.T) {
	var cases = map[string]struct {
		session Session
		wantErr bool
	}{
		"Both are empty": {
			Session{},
			false,
		},
		"Both contain valid values": {
			Session{
				Analysis: "sv",
				List:     "aml",
			},
			false,
		},
		"Analysis is empty": {
			Session{
				Analysis: "",
				List:     "all",
			},
			true,
		},
		"List is empty": {
			Session{
				Analysis: "sv",
				List:     "",
			},
			true,
		},
		"Analysis contains invalid value": {
			Session{
				Analysis: "stv",
				List:     "aml",
			},
			true,
		},
		"List contains invalid value": {
			Session{
				Analysis: "snv",
				List:     "atl",
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.session.checkListAndAnalysis()
			checkError(t, err, c.wantErr)
		})
	}
}

func TestModifyFlagInput(t *testing.T) {
	var cases = map[string]struct {
		session Session
		result  Session
	}{
		"Bed and List field contain values": {
			Session{
				Analysis: "snv",
				Bed:      "my.bed",
				List:     "aml",
			},
			Session{
				Analysis: "snv",
				Bed:      "my.bed",
				DbList:   "list_aml",
				List:     "aml",
			},
		},
		"Bed field is empty": {
			Session{
				Analysis: "snv",
				Bed:      "",
				List:     "aml",
			},
			Session{
				Analysis: "snv",
				Bed:      "aml_snv.bed",
				List:     "aml",
				DbList:   "list_aml",
			},
		},
		"List field is empty": {
			Session{
				Analysis: "snv",
				Bed:      "my.bed",
				List:     "",
			},
			Session{
				Analysis: "snv",
				Bed:      "my.bed",
				List:     "",
				DbList:   "",
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.session.modifyFlagInput()
			if diff := deep.Equal(c.session, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestReadFlags(t *testing.T) {
	var cases = map[string]struct {
		input   []string
		result  Session
		wantErr bool
	}{
		"All flags set": {
			[]string{"", "-analysis", "snv", "-bed", "aml_snv.bed", "-list", "aml", "-tsv", "aml_snv.tsv"},
			Session{
				Analysis: "snv",
				Bed:      "aml_snv.bed",
				DbList:   "list_aml",
				List:     "aml",
				Tsv:      "aml_snv.tsv",
			},
			false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var session Session
			os.Args = c.input
			flag.CommandLine = flag.NewFlagSet("Reset", flag.ExitOnError)
			err := session.readFlags()
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(session, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}
