package license

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/report"
)

var (
	ArgsError         = errors.New("failed to handle args")
	SubscriptionError = errors.New("premium feature. Please visit https://debricked.com/pricing/ for more info")
)

type OrderArgs struct {
	Email      string
	CommitHash string
}

type Reporter struct {
	DebClient client.IDebClient
}

func (r Reporter) Order(args report.IOrderArgs) error {
	orderArgs, ok := args.(OrderArgs)
	if !ok {
		return ArgsError
	}

	commitId, err := r.getCommitId(orderArgs.CommitHash)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("/api/1.0/open/licenses/get-licenses?order=asc&sortColumn=name&generateExcel=1&commitId=%d&email=%s", commitId, orderArgs.Email)
	res, err := r.DebClient.Get(uri, "application/json")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusForbidden {
		return SubscriptionError
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to order report. ReceiveStatus code: %d", res.StatusCode)
	}

	return nil
}

type commit struct {
	FileIds     []int  `json:"uploaded_programs_file_ids"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ReleaseData string `json:"release_date"`
}

func (r Reporter) getCommitId(hash string) (int, error) {
	uri := fmt.Sprintf("/api/1.0/releases/by/name?name=%s", hash)
	res, err := r.DebClient.Get(uri, "application/json")
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return 0, SubscriptionError
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("no commit was found with the name %s", hash)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var commits []commit
	err = json.Unmarshal(body, &commits)
	if len(commits) == 0 {
		return 0, fmt.Errorf("no commit was found with the name %s", hash)
	}

	return commits[0].Id, err
}
