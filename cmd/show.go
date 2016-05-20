package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var raw bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <number>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You need to ask for one issue by number!")
		}
		issue, err := db.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if raw {
			b, err := json.MarshalIndent(issue, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(b))
		} else {
			fmt.Print(issue.FmtTitle())
			fmt.Print(issue.FmtByLine())
			if issue.Body != nil {
				fmt.Printf("\n%s\n", *issue.Body)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Show the raw JSON for this issue")
}
