package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (s Session) tsvToDb() (err error) {
	tsv, err := readTsv(s.Tsv)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Could not read file %s", s.Tsv))
		return
	}
	var header []string
	for i, row := range tsv {
		if i == 0 {
			if err = validateTsvHeader(row); err != nil {
				return
			}
		} else {
			err = addRowToDb(row, header)
		}
		//for _, table := range dbRow.Tables {
		//	if err = s.ensureTableExists(table); err != nil {
		//		return
		//	}
		//	if err = s.checkAndAddEntity(table, entity); err != nil {
		//		return
		//	}
		//	if err = s.prepAndAddFusions(table, entity); err != nil {
		//		return
		//	}
	}
	return
}

func validateTsvHeader(header []string) (err error) {
	for _, column := range header {
		if _, present := tsvHeader[column]; !present {
			err = errors.New(fmt.Sprintf("Found unknown column: %s", column))
			return
		}
	}
	var missing []string
	for key, required := range tsvHeader {
		for i, column := range header {
			if !required {
				break
			} else if key == column {
				break
			} else if len(header) == i+1 {
				missing = append(missing, key)
			}
		}
	}
	if len(missing) > 0 {
		err = errors.New(fmt.Sprintf("The following columns are missing: %s", strings.Join(missing, ", ")))
	}
	return
}

func addRowToDb(row []string, header []string) (err error) {
	var mpRow map[string]string
	mpRow, err = rowToMap(row, header)
	if err != nil {
		return
	}
	var dbRow DbTableRow
	dbRow, err = mapToDbRow(mpRow)
	if err != nil {
		return
	}
	for _, table := range dbRow.Tables {
		table = fmt.Sprintf("table_%s", strings.ToLower(table))
	}
	return
}

func rowToMap(row []string, header []string) (result map[string]string, err error) {
	if len(row) != len(header) {
		err = errors.New(fmt.Sprintf("Header and row length are differing for row: %s", strings.Join(row, " ")))
		return
	}
	result = make(map[string]string)
	for i, column := range header {
		result[column] = row[i]
	}
	return
}

func mapToDbRow(row map[string]string) (dbRow DbTableRow, err error) {
	if err = dbRow.validateAnalyses(strings.ToLower(row["analyses"])); err != nil {
		return
	}
	if err = dbRow.validateClass(strings.ToLower(row["class"])); err != nil {
		return
	}
	if err = dbRow.validateCoordinates(row["coordinates"]); err != nil {
		return
	}
	dbRow.Id = row["id"]
	if err = dbRow.validateIncludePartners(row["include_partners"]); err != nil {
		return
	}
	dbRow.Tables = strings.Split(row["tables"], ",")
	return
}

func (d *DbTableRow) validateAnalyses(anString string) (err error) {
	d.Analyses = make(map[string]struct{})
	for _, analysis := range strings.Split(anString, ",") {
		if _, valid := analyses[analysis]; !valid {
			err = errors.New(fmt.Sprintf("%s is not a valid analysis", analysis))
			break
		} else {
			d.Analyses[analysis] = struct{}{}
		}
	}
	return
}

func (d *DbTableRow) validateClass(class string) (err error) {
	if _, valid := classes[class]; !valid {
		err = errors.New(fmt.Sprintf("%s is not a valid class", class))
	} else {
		d.Class = class
	}
	return
}

func (d *DbTableRow) validateCoordinates(coordinates string) (err error) {
	if coordinates == "" {
		return
	}
	coordRegex := regexp.MustCompile("^(chr)?[\\d,X,Y,M]\\d?:\\d+-\\d+$")
	if !coordRegex.MatchString(coordinates) {
		err = errors.New(fmt.Sprintf("%s does not match expected coordinates string (e.g. chr1:0-10)", coordinates))
		return
	}
	coordSlice := strings.FieldsFunc(coordinates, split)
	chromosome := strings.ReplaceAll(coordSlice[0], "chr", "")
	if err = d.validateChromosome(chromosome); err != nil {
		return
	}
	d.Start = coordSlice[1]
	d.End = coordSlice[2]
	return
}

func split(r rune) bool {
	return r == ':' || r == '-'
}

func (d *DbTableRow) validateChromosome(chromosome string) (err error) {
	chromosomes := generateChromosomeMap()
	if !chromosomes[chromosome] {
		err = errors.New(fmt.Sprintf("%s is not a valid chromosome", chromosome))
		return
	}
	d.Chromosome = chromosome
	return
}

func generateChromosomeMap() (chromosomes map[string]bool) {
	chromosomes = make(map[string]bool)
	for i := 1; i <= 22; i++ {
		chromosomes[strconv.Itoa(i)] = true
	}
	for _, chromosome := range []string{"M", "X", "Y"} {
		chromosomes[chromosome] = true
	}
	return
}

func (d *DbTableRow) validateIncludePartners(include string) (err error) {
	d.IncludePartners, err = strconv.ParseBool(include)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("%s could not be converted to a valid bool", include))
	}
	return
}
