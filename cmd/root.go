package cmd

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

var session Session

func init() {
	if err := env.Parse(&session); err != nil {
		log.Fatalf("%v", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gene_list_svc",
	Short: "Update gene lists and generate bed files",
	Long:  `Update database with genetic regions in tsv files and generate bed files from specific tables`,
}
