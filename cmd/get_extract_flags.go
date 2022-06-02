package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func getExtractFlags(cmd cobra.Command) (err error) {
	if err = validateAnalysis(cmd); err != nil {
		return
	}
	if err = validateBuild(cmd); err != nil {
		return
	}
	if err = validateTables(cmd); err != nil {
		return
	}
	if err = getBedName(cmd); err != nil {
		return
	}
	if session.Chr, err = cmd.Flags().GetBool("chr"); err != nil {
		return
	}
	return
}

func validateAnalysis(cmd cobra.Command) (err error) {
	analysis, err := cmd.Flags().GetString("analysis")
	if err != nil {
		return
	}
	if _, valid := analyses[analysis]; !valid {
		err = errors.New(fmt.Sprintf("%s is not a valid analysis", analysis))
		return
	} else {
		session.Analysis = analysis
	}
	return
}

func validateBuild(cmd cobra.Command) (err error) {
	build, err := cmd.Flags().GetString("build")
	if err != nil {
		return
	}
	if build == "38" || build == "37" {
		session.Build = build
	} else {
		err = errors.New(fmt.Sprintf("%s is not a valid genome build", build))
	}
	return
}

func validateTables(cmd cobra.Command) (err error) {
	dbTables, err := session.Db.Connection.getTables()
	if err != nil {
		return
	}
	tables, err := cmd.Flags().GetString("tables")
	for _, table := range strings.Split(tables, ",") {
		if _, valid := dbTables[table]; !valid {
			err = errors.New(fmt.Sprintf("table %s is not present in database", table))
			return
		} else {
			session.Tables = append(session.Tables, table)
		}
	}
	return
}

func getBedName(cmd cobra.Command) (err error) {
	bed, err := cmd.Flags().GetString("bed")
	if err != nil {
		return
	}
	if bed != "" {
		session.Bed = bed
	} else {
		session.Bed = fmt.Sprintf("%s_%s_%s_%s.bed", strings.Join(session.Tables, "_"), session.Analysis, session.Build, time.Now().Format("2006-01-02"))
	}
	return
}
