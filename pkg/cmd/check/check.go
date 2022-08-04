package check

import (
	"debricked/pkg/client"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
)

var debClient client.Client

func NewCheckCmd(debrickedClient *client.Client) *cobra.Command {
	debClient = *debrickedClient
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

type latestScan struct {
	Id   *int    `json:"id"`
	Date *int    `json:"date"`
	Url  *string `json:"commitUrl"`
}

func getScanId(commitId int) (int, error) {
	uri := fmt.Sprintf("/api/1.0/open/scan/latest-scan-status?commitId=%d", commitId)
	res, err := debClient.Get(uri, "application/json")
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("No scan was found for commit"))
	}

	body, err := io.ReadAll(res.Body)
	var scan map[string]latestScan
	err = json.Unmarshal(body, &scan)
	if err != nil {
		return 0, err
	}
	id := scan["latestScan"].Id
	if id == nil {
		return 0, errors.New(fmt.Sprintf("No scan was found for commit"))
	}

	return *id, err
}
