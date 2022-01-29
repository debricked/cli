package check

import (
	"debricked/pkg/client"
	"github.com/spf13/cobra"
	"log"
)

var debClient *client.DebClient

func NewCheckCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	cmd := &cobra.Command{
		Use:   "check [commit hash]",
		Short: "Check scan results on a specific commit",
		Long: `Check fetch identified vulnerabilities and Automations results from a
specific commit`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			commitHash := args[0]
			err := check(commitHash)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}

func check(hash string) error {

	return nil
}
