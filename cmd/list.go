package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/markbates/ghi/cmd/issue"
	"github.com/spf13/cobra"
)

var state string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists issues for the repo.",
	Run: func(cmd *cobra.Command, args []string) {
		var issues []issue.Issue
		var err error
		switch state {
		case "all":
			issues, err = db.All()
		default:
			issues, err = db.AllByState(state)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if raw {
			b, err := json.MarshalIndent(issues, "", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Print(string(b))
		} else {
			for _, issue := range issues {
				fmt.Print(issue.FmtTitle())
			}
			fmt.Printf("\n=== (%d) Issues ===\n", len(issues))
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Show the raw JSON for these issues")
	listCmd.Flags().StringVarP(&state, "state", "s", "open", "List issues by their state <all, closed, open>")
}
