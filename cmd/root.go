package cmd

import (
	"fmt"
	"os"

	"github.com/markbates/ghi/cmd/store"
	"github.com/spf13/cobra"
)

var db *store.Store
var config *Config

var RootCmd = &cobra.Command{
	Use:   "ghi",
	Short: "Offline GitHub Issues Client",
	Long:  `GHI let's you download issues from a GitHub repo to be made available offline.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config = LoadConfig()
		if config.Repo == "" {
			if len(args) == 0 {
				fmt.Println(`It looks like you haven't initialized GHI yet!

The first time you run GHI you should run "ghi fetch owner/repo".

This will fetch all of your issues for that repository. Future calls
to "fetch" won't require the "owner/repo" since we'll store a little
meta-data file in this repo to track that.`)
				os.Exit(-1)
			} else {
				config.SetFromArgs(args)
			}
		}

		s, err := store.New(config.Owner, config.Repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		db = s
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
