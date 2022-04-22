package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Add genetic regions to database",
	Long:  `Add genetic regions, specified in a tsv file, to corresponding list in database`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		session.Tsv, err = cmd.Flags().GetString("tsv")
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func init() {
	// Add update command
	rootCmd.AddCommand(updateCmd)

	// Add tsv flag to update command
	updateCmd.PersistentFlags().String("tsv", "", "tsv containg list of genetic regions")
}
