package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/go-github/github"
	"github.com/markbates/ghi/cmd/issue"
	"github.com/markbates/going/wait"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			config.SetFromArgs(args)
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)
		opts := &github.IssueListByRepoOptions{State: "all"}
		allIssues := []issue.Issue{}
		more := true
		locker := &sync.Mutex{}
		for more {
			issues, resp, err := client.Issues.ListByRepo(db.Owner, db.Repo, opts)
			if err != nil {
				log.Fatal(err)
			}

			wait.Wait(len(issues), func(i int) {
				issue := &issue.Issue{Issue: issues[i], Comments: []github.IssueComment{}}
				comments, _, err := client.Issues.ListComments(db.Owner, db.Repo, *issue.Number, &github.IssueListCommentsOptions{})
				if err != nil {
					log.Fatal(err)
				}
				issue.Comments = comments
				locker.Lock()
				allIssues = append(allIssues, *issue)
				locker.Unlock()
			})
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		err := db.Persist(allIssues)
		if err != nil {
			log.Fatal(err)
		}
		config.Save()
		fmt.Printf("Fetched %d issues for %s/%s\n", len(allIssues), db.Owner, db.Repo)
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)
}
