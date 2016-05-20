package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		issues, err := db.All()
		if err != nil {
			log.Fatal(err)
		}
		if raw {
			b, err := json.MarshalIndent(issues, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(b))
		} else {
			for _, issue := range issues {
				fmt.Print(issue.FmtTitle())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Show the raw JSON for these issues")
}
