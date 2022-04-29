package cmd

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract genetic regions from database",
	Long:  `Extract genetic regions from corresponding lists in database and generate bed file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := session.initDbConnection(); err != nil {
			log.Fatalf("%v", err)
		}
		if err := getExtractFlags(*cmd); err != nil {
			log.Fatalf("%v", err)
		}
		if err := dbToTsv(); err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func init() {
	// Add extract command
	rootCmd.AddCommand(extractCmd)

	// Add flags to extract command
	extractCmd.PersistentFlags().String("analysis", "", "choose analysis (cnv, pindel, snv, sv)")
	extractCmd.PersistentFlags().String("bed", "", `set individual bed file name (default "tables_analysis_build_timestamp.bed")`)
	extractCmd.PersistentFlags().String("build", "38", "choose genome build")
	extractCmd.PersistentFlags().Bool("chr", true, "use chr-prefix for chromosome ids")
	extractCmd.PersistentFlags().String("tables", "", "comma-separated list of tables to be included")
}
