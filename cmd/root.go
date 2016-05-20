package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/ghi/cmd/store"
	"github.com/spf13/cobra"
)

var db *store.Store
var config *Config

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ghi",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config = LoadConfig()
		if config.Repo == "" {
			if len(args) == 0 {
				log.Fatal("You must pass in an owner/repo to the this function!")
			} else {
				config.SetFromArgs(args)
			}
		}

		s, err := store.New(config.Owner, config.Repo)
		if err != nil {
			log.Fatal(err)
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
