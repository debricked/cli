package license

import (
	"debricked/pkg/client"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

var debClient *client.DebClient

var email string
var commitHash string

func NewLicenseCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Generate license report",
		Long: `Generate license report from a commit hash. 
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.
The finished report will be sent to the specified email address.`,
		RunE: run,
	}

	cmd.Flags().StringVarP(&email, "email", "e", "", "The email address that the report will be sent to")
	_ = cmd.MarkFlagRequired("email")

	cmd.Flags().StringVarP(&commitHash, "commit", "c", "", "commit hash")
	_ = cmd.MarkFlagRequired("commit")

	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	if err := report(); err != nil {
		return errors.New(fmt.Sprintf("%s %s\n", color.RedString("⨯"), err.Error()))
	}

	fmt.Printf("%s Successfully ordered license report\n", color.GreenString("✔"))

	return nil
}

func report() error {
	commitId, err := getCommitId(commitHash)
	if err != nil {
		return err
	}

	return orderReport(commitId)
}

func orderReport(commitId int) error {
	uri := fmt.Sprintf("/api/1.0/open/licenses/get-licenses?order=asc&sortColumn=name&generateExcel=1&commitId=%d&email=%s", commitId, email)
	res, err := debClient.Get(uri, "application/json")
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusForbidden {
		return errors.New("premium feature. Please visit https://debricked.com/pricing/ for more info")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to order report. Status code: %d", res.StatusCode))
	}

	return nil
}

type commit struct {
	FileIds     []int  `json:"uploaded_programs_file_ids"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ReleaseData string `json:"release_date"`
}

func getCommitId(hash string) (int, error) {
	uri := fmt.Sprintf("/api/1.0/releases/by/name?name=%s", hash)
	res, err := debClient.Get(uri, "application/json")
	if err != nil {
		return 0, err
	}

	if res.StatusCode == http.StatusForbidden {
		return 0, errors.New("premium feature. Please visit https://debricked.com/pricing/ for more info")
	}

	if res.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("No commit was found with the name %s", hash))
	}

	body, err := io.ReadAll(res.Body)
	var commits []commit
	err = json.Unmarshal(body, &commits)
	if len(commits) == 0 {
		return 0, errors.New(fmt.Sprintf("No commit was found with the name %s", hash))
	}

	return commits[0].Id, err
}
