package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/markbates/ghi/cmd/issue"
	"github.com/markbates/going/wait"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var fetchState string

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches all of the issues for the specified repo.",
	Long: `Fetches all of the issues for the specified repo.
This will clear any existing issues stored locally, pull down
everything from GitHub and store it locally for offline use.

The first time you run this command you should run it like such:

$ ghi fetch owner/repo

Subsequent calls will not need the "owner/repo".

If you are going to be calling a private repo you will need to
set the ENV var "GITHUB_TOKEN" with a GitHub Personal Access Token.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			config.SetFromArgs(args)
		}

		err := db.Clear()
		if err != nil {
			log.Fatal(err)
		}

		more := true
		count := 0
		opts := &github.IssueListByRepoOptions{State: fetchState}

		spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		spin.Suffix = fmt.Sprintf(" Fetching %s Issues (This could take a while depending on the number of issues and comments you have!)", fetchState)
		spin.Start()

		client := newClient()
		for more {
			issues, resp, err := client.Issues.ListByRepo(db.Owner, db.Repo, opts)
			if err != nil {
				fmt.Println(err)
				if resp != nil && resp.StatusCode == 401 {
					fmt.Println(`Couldn't access this repo! Try setting a GitHub Personal Token.

This token can be set as an environment variable "GITHUB_TOKEN".`)
				}
				os.Exit(-1)
			}
			count += len(issues)

			wait.Wait(len(issues), func(i int) {
				issue := &issue.Issue{Issue: issues[i], Comments: []github.IssueComment{}}
				comments, _, err := client.Issues.ListComments(db.Owner, db.Repo, *issue.Number, &github.IssueListCommentsOptions{})
				if err != nil {
					log.Fatal(err)
				}
				issue.Comments = comments
				db.Save(*issue)
			})
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		if err != nil {
			log.Fatal(err)
		}
		config.Save()
		spin.Stop()
		fmt.Printf("\nFetched %d issues for %s/%s\n", count, db.Owner, db.Repo)
	},
}

func newClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&fetchState, "state", "s", "all", "Fetch issues by their state <all, closed, open>")
}
