package sbom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/debricked/cli/internal/client"
	internalIO "github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/report"
	"github.com/fatih/color"
)

var (
	ErrHandleArgs   = errors.New("failed to handle args")
	ErrSubscription = errors.New("premium feature. Please visit https://debricked.com/pricing/ for more info")
)

type generateSbom struct {
	Format                string   `json:"format"`
	RepositoryID          string   `json:"repositoryId"`
	IntegrationName       string   `json:"integrationName"`
	CommitID              string   `json:"commitId"`
	Email                 string   `json:"email"`
	Branch                string   `json:"branch"`
	Locale                string   `json:"locale"`
	Licenses              bool     `json:"licenses"`
	Vulnerabilities       bool     `json:"vulnerabilities"`
	SendEmail             bool     `json:"sendEmail"`
	VulnerabilityStatuses []string `json:"vulnerabilityStatuses"`
}

type generateSbomResponse struct {
	Message    string   `json:"message"`
	ReportUUID string   `json:"reportUuid"`
	Notes      []string `json:"notes"`
}

type OrderArgs struct {
	RepositoryID    string
	CommitID        string
	Branch          string
	Output          string
	Format          string
	Vulnerabilities bool
	Licenses        bool
}

type Reporter struct {
	DebClient  client.IDebClient
	FileWriter internalIO.IFileWriter
}

func (r Reporter) Order(args report.IOrderArgs) error {
	orderArgs, ok := args.(OrderArgs)
	var err error
	if !ok {
		return ErrHandleArgs
	}

	uuid, err := r.generate(orderArgs)
	if err != nil {
		return err
	}
	sbom, err := r.download(uuid)
	if err != nil {
		return err
	}

	return r.writeSBOM(orderArgs.Output, orderArgs.RepositoryID, orderArgs.CommitID, sbom)

}

func (r Reporter) generate(orderArgs OrderArgs) (string, error) {
	// Tries to start generating an SBOM and returns the UUID for the report
	body, err := json.Marshal(generateSbom{
		Format:                orderArgs.Format,
		RepositoryID:          orderArgs.RepositoryID,
		CommitID:              orderArgs.CommitID,
		Email:                 "",
		Branch:                orderArgs.Branch,
		Locale:                "en",
		Vulnerabilities:       orderArgs.Vulnerabilities,
		Licenses:              orderArgs.Licenses,
		SendEmail:             false,
		VulnerabilityStatuses: []string{"vulnerable", "unexamined", "paused", "snoozed"},
	})

	if err != nil {
		return "", err
	}

	response, err := (r.DebClient).Post(
		"/api/1.0/open/sbom/generate",
		"application/json",
		bytes.NewBuffer(body),
		0,
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusPaymentRequired {
		return "", ErrSubscription
	} else if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to initialize SBOM generation due to status code %d", response.StatusCode)
	} else {
		fmt.Println("Successfully initialized SBOM generation")
	}

	return r.parseUUID(response.Body)
}

func (r Reporter) parseUUID(body io.Reader) (string, error) {
	generateSbomResponseJSON, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	var generateSbomResponse generateSbomResponse
	err = json.Unmarshal(generateSbomResponseJSON, &generateSbomResponse)

	return generateSbomResponse.ReportUUID, err
}

func (r Reporter) download(uuid string) ([]byte, error) {
	uri := fmt.Sprintf("/api/1.0/open/sbom/download?reportUuid=%s", uuid)
	fmt.Printf("%s", color.BlueString("Downloading SBOM..."))
	for { // poll download status until completion
		res, err := (r.DebClient).Get(uri, "application/json")

		if err != nil {
			return nil, err
		}
		switch statusCode := res.StatusCode; statusCode {
		case http.StatusOK:
			data, _ := io.ReadAll(res.Body)
			defer res.Body.Close()
			fmt.Printf("%s\n", color.GreenString("✔"))

			return data, nil
		case http.StatusCreated:
			return nil, errors.New("polling failed due to too long queue times")
		case http.StatusAccepted:
			time.Sleep(5000 * time.Millisecond)
		default:
			return nil, fmt.Errorf("download failed with status code %d", res.StatusCode)
		}
	}
}

func (reporter Reporter) writeSBOM(output, repositoryID, commitID string, sbomBytes []byte) error {
	var filename string
	if output == "" {
		filename = fmt.Sprintf("%s-%s.sbom.json", repositoryID, commitID)
	} else {
		filename = output
	}
	file, err := reporter.FileWriter.Create(filename)
	if err != nil {
		return err
	}

	return reporter.FileWriter.Write(file, sbomBytes)
}

func (reporter Reporter) ParseDetailsURL(detailsURL string) (string, string, error) {
	// Parses CommitID and RepositoryID from the details URL which has the format;
	// https://debricked.com/app/en/repository/<repository_id>/commit/<commit_id>"
	urlParts := strings.Split(detailsURL, "/")
	if len(urlParts) != 9 {

		return "", "", fmt.Errorf("URL \"%s\"is of wrong format", detailsURL)
	}

	return urlParts[6], urlParts[8], nil
}
