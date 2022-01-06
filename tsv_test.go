package main

import (
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
)

func TestReadTsv(t *testing.T) {
	var cases = map[string]struct {
		path    string
		result  [][]string
		wantErr bool
	}{
		"File exists": {
			"test_data/test.tsv",
			[][]string{
				{
					"ABL1",
					"Gene",
					"snv",
					"aml",
					"Specific",
				},
			},
			false,
		},
		"File does not exist": {
			"test_data/not_existent.tsv",
			nil,
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := readTsv(c.path)
			checkError(t, err, c.wantErr)
			if diff := deep.Equal(result, c.result); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestWriteTsv(t *testing.T) {
	var cases = map[string]struct {
		path    string
		data    [][]string
		wantErr bool
	}{
		"Write data successfully": {
			"test_data/test_output.tsv",
			[][]string{
				{
					"ABL1",
					"Gene",
					"snv",
					"aml",
					"Specific",
				},
			},
			false,
		},
		"Data writing fails": {
			"test_data/not_existent/test_output.tsv",
			[][]string{
				{
					"ABL1",
					"Gene",
					"snv",
					"aml",
					"Specific",
				},
			},
			true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := writeTsv(c.path, c.data)
			checkError(t, err, c.wantErr)
			if _, err = os.Stat(c.path); err == nil {
				if err = os.Remove(c.path); err != nil {
					log.Fatalf("%v", err)
				}
			}
		})
	}
}
