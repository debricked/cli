package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/git"
	"github.com/debricked/cli/internal/tui"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var (
	NoFilesErr           = errors.New("failed to find dependency files")
	PollingTerminatedErr = errors.New("progress polling terminated due to long queue times")
	EmptyFileErr         = errors.New("tried to upload empty file")
	InitScanErr          = errors.New("failed to initialize a scan")
)

const callgraphName = "debricked-call-graph"

type uploadBatch struct {
	client             *client.IDebClient
	fileGroups         file.Groups
	gitMetaObject      *git.MetaObject
	integrationName    string
	ciUploadId         int
	callGraphTimeout   int
	versionHint        bool
	debrickedConfig    *DebrickedConfig // JSON Config
	tagCommitAsRelease bool
	experimental       bool
}

func newUploadBatch(
	client *client.IDebClient, fileGroups file.Groups, gitMetaObject *git.MetaObject,
	integrationName string, callGraphTimeout int, versionHint bool,
	debrickedConfig *DebrickedConfig, tagCommitAsRelease bool, experimental bool,
) *uploadBatch {
	return &uploadBatch{
		client:             client,
		fileGroups:         fileGroups,
		gitMetaObject:      gitMetaObject,
		integrationName:    integrationName,
		ciUploadId:         0,
		callGraphTimeout:   callGraphTimeout,
		versionHint:        versionHint,
		debrickedConfig:    debrickedConfig,
		tagCommitAsRelease: tagCommitAsRelease,
		experimental:       experimental,
	}
}

// upload concurrently posts all file groups to Debricked
func (uploadBatch *uploadBatch) upload() error {
	uploadWorker := func(fileQueue <-chan string, fileResults chan<- int) {
		const ok = 0
		const fail = 1

		for f := range fileQueue {
			fileName := filepath.Base(f)
			var err error
			timeout := 0
			if strings.HasSuffix(fileName, callgraphName) {
				timeout = uploadBatch.callGraphTimeout
			}
			err = uploadBatch.uploadFile(f, timeout)

			if err != nil {
				log.Println("Failed to upload:", f)
				if err != nil {
					log.Println(err.Error())
					fileResults <- fail
				}
			} else {
				printSuccessfulUpload(f)
				fileResults <- ok
			}
		}
	}

	files, err := uploadBatch.initUpload()
	if err != nil {

		return err
	}

	fileQueue := make(chan string, len(files))
	fileResults := make(chan int, len(files))

	// Spawn workers
	for w := 1; w <= 20; w++ {
		go uploadWorker(fileQueue, fileResults)
	}

	// Append file jobs on queue
	for _, f := range files {
		fileQueue <- f
	}

	// Await completion
	for range files {
		<-fileResults
	}

	close(fileQueue)

	return nil
}

// uploadFile Reads file content from filepath and uploads it to Debricked. Returns HTTP status code or 0 if other error occur
func (uploadBatch *uploadBatch) uploadFile(filePath string, timeout int) error {
	if strings.HasSuffix(filePath, "debricked.fingerprints.txt") && !(*uploadBatch.client).IsEnterpriseCustomer(true) {
		return errors.New("non-enterprise customer trying to upload fingerprints")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileData, _ := writer.CreateFormFile("fileData", filepath.Base(filePath))
	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, _ = io.Copy(fileData, f)

	_ = writer.WriteField("fileRelativePath", getRelativeFilePath(filePath))
	_ = writer.WriteField("repositoryName", uploadBatch.gitMetaObject.RepositoryName)
	_ = writer.WriteField("commitName", uploadBatch.gitMetaObject.CommitName)
	_ = writer.WriteField("repositoryUrl", uploadBatch.gitMetaObject.RepositoryUrl)
	_ = writer.WriteField("branchName", uploadBatch.gitMetaObject.BranchName)
	if uploadBatch.initialized() {
		_ = writer.WriteField("ciUploadId", strconv.Itoa(uploadBatch.ciUploadId))
	}
	response, err := (*uploadBatch.client).Post(
		"/api/1.0/open/uploads/dependencies/files",
		writer.FormDataContentType(),
		body,
		timeout,
	)
	if err != nil {
		return err
	}
	if !uploadBatch.initialized() {
		data, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		uFile := uploadedFile{}
		_ = json.Unmarshal(data, &uFile)
		if uFile.CiUploadId == 0 {
			return EmptyFileErr
		}
		uploadBatch.ciUploadId = uFile.CiUploadId
	}

	return nil
}

// initAnalysis send the finish request that starts the analysis
func (uploadBatch *uploadBatch) initAnalysis() error {
	if uploadBatch.ciUploadId == 0 {
		return NoFilesErr
	}
	body, err := json.Marshal(uploadFinish{
		CiUploadId:           strconv.Itoa(uploadBatch.ciUploadId),
		RepositoryName:       uploadBatch.gitMetaObject.RepositoryName,
		IntegrationName:      uploadBatch.integrationName,
		CommitName:           uploadBatch.gitMetaObject.CommitName,
		Author:               uploadBatch.gitMetaObject.Author,
		VersionHint:          uploadBatch.versionHint,
		DebrickedConfig:      uploadBatch.debrickedConfig,
		DebrickedIntegration: "cli",
		TagCommitAsRelease:   uploadBatch.tagCommitAsRelease,
		Experimental:         uploadBatch.experimental,
	})

	if err != nil {
		return err
	}

	response, err := (*uploadBatch.client).Post(
		"/api/1.0/open/finishes/dependencies/files/uploads",
		"application/json",
		bytes.NewBuffer(body),
		0,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to initialize scan due to status code %d", response.StatusCode)
	} else {
		fmt.Println("Successfully initialized scan")
	}

	return nil
}

func (uploadBatch *uploadBatch) initialized() bool {
	return uploadBatch.ciUploadId > 0
}

// wait track scan progress and return uploadStatus upon completion
func (uploadBatch *uploadBatch) wait() (*UploadResult, error) {
	bar := tui.NewProgressBar()
	_ = bar.RenderBlank()
	// poll scan status until completion
	var resultStatus *UploadResult
	uri := fmt.Sprintf("/api/1.0/open/ci/upload/status?ciUploadId=%s", strconv.Itoa(uploadBatch.ciUploadId))
	for !bar.IsFinished() {
		res, err := (*uploadBatch.client).Get(uri, "application/json")
		if err != nil {
			return nil, err
		}
		status, err := newUploadStatus(res)
		if err != nil {
			return nil, err
		}
		if res.StatusCode == http.StatusCreated {
			err := bar.Finish()
			if err != nil {
				return resultStatus, err
			}

			resultStatus = &UploadResult{
				DetailsUrl: status.DetailsUrl,
				LongQueue:  true,
			}

			return resultStatus, PollingTerminatedErr
		}
		err = bar.Set(status.Progress)
		if err != nil {
			return nil, err
		}

		if bar.IsFinished() {
			resultStatus = newUploadResult(status)
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}

	return resultStatus, nil
}

// initUpload initialises a scan by uploading one file. This enables the scan to
// get assigned a `ciUploadId`
func (uploadBatch *uploadBatch) initUpload() ([]string, error) {
	files := uploadBatch.fileGroups.GetFiles()
	if len(files) == 0 {
		return files, nil
	}

	var entryFile string
	var err error
	for len(files) > 0 {
		entryFile = files[0]
		files = files[1:]
		timeout := 0
		if strings.HasSuffix(filepath.Base(entryFile), callgraphName) {
			timeout = 30
		}
		err = uploadBatch.uploadFile(entryFile, timeout)
		if err == nil {
			printSuccessfulUpload(entryFile)

			return files, nil
		}
	}

	errStr := fmt.Sprintf("Failed to initialize a scan for %s. Got the following error: %s", entryFile, err.Error())

	return files, errors.New(errStr)
}

type uploadedFile struct {
	CiUploadId           int    `json:"ciUploadId"`
	UploadProgramsFileId int    `json:"uploadProgramsFileId"`
	TotalScans           int    `json:"totalScans"`
	RemainingScans       int    `json:"remainingScans"`
	Percentage           string `json:"percentage"`
	EstimateDaysLeft     int    `json:"estimateDaysLeft"`
}

type boolOrString struct {
	Version    string `json:"version"`
	HasVersion bool   `json:"hasVersion"`
}

func (boolOrString *boolOrString) MarshalJSON() ([]byte, error) {
	if !boolOrString.HasVersion {
		return json.Marshal(&boolOrString.HasVersion)
	}

	return json.Marshal(&boolOrString.Version)
}

type purlConfig struct {
	PackageURL  string       `json:"pURL" yaml:"pURL"`
	Version     boolOrString `json:"version" yaml:"version"` // Either false or version string
	FileRegexes []string     `json:"fileRegexes" yaml:"fileRegexes"`
}

type DebrickedConfig struct {
	Overrides []purlConfig  `json:"override,omitempty" yaml:"overrides"`
	Ignore    *IgnoreConfig `json:"ignore,omitempty" yaml:"ignore,omitempty"`
}

// IgnoreConfig matches the structure of the 'ignore' section in YAML
type IgnoreConfig struct {
	Packages []IgnorePackage `json:"packages" yaml:"packages"`
}

type IgnorePackage struct {
	PURL    string `json:"pURL" yaml:"pURL"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

type uploadFinish struct {
	CiUploadId           string           `json:"ciUploadId"`
	RepositoryName       string           `json:"repositoryName"`
	IntegrationName      string           `json:"integrationName"`
	CommitName           string           `json:"commitName"`
	Author               string           `json:"author"`
	DebrickedIntegration string           `json:"debrickedIntegration"`
	VersionHint          bool             `json:"versionHint"`
	DebrickedConfig      *DebrickedConfig `json:"debrickedConfig"`
	TagCommitAsRelease   bool             `json:"isRelease"`
	Experimental         bool             `json:"experimental"`
}

func getRelativeFilePath(filePath string) string {
	relFilePath := filepath.Dir(filePath)
	if strings.EqualFold(".", relFilePath) {
		relFilePath = ""
	}

	return relFilePath
}

func printSuccessfulUpload(f string) {
	fmt.Printf("Successfully uploaded: %s\n", color.YellowString(f))
}

type pURLConfigYAML struct {
	PackageURL  string   `yaml:"pURL"`
	Version     *string  `yaml:"version"`
	FileRegexes []string `yaml:"fileRegexes"`
}

type DebrickedConfigYAML struct {
	Overrides []pURLConfigYAML `yaml:"overrides"`
}

// extractIgnore unmarshals the ignore section from raw config
func extractIgnore(raw map[string]interface{}) *IgnoreConfig {
	if rawIgnore, ok := raw["ignore"]; ok {
		ignoreYaml, err := yaml.Marshal(rawIgnore)
		if err == nil {
			var ignoreObj IgnoreConfig
			if yaml.Unmarshal(ignoreYaml, &ignoreObj) == nil {
				return &ignoreObj
			}
		}
	}

	return nil
}

// convertOverrides converts YAML overrides to purlConfig slice
func convertOverrides(yamlOverrides []pURLConfigYAML) []purlConfig {
	var overrides []purlConfig
	for _, entry := range yamlOverrides {
		var version string
		var exist bool
		pURL := entry.PackageURL
		fileRegexes := entry.FileRegexes
		if entry.Version == nil {
			version = ""
			exist = false
		} else {
			version = *entry.Version
			exist = true
		}
		overrides = append(overrides, purlConfig{PackageURL: pURL, Version: boolOrString{Version: version, HasVersion: exist}, FileRegexes: fileRegexes})
	}

	return overrides
}

func GetDebrickedConfig(path string) *DebrickedConfig {
	var yamlConfig DebrickedConfigYAML
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf(
			"%s Failed to read debricked config file on path \"%s\"",
			color.YellowString("⚠️"),
			path,
		)

		return nil
	}

	// Unmarshal into map to support any key for overrides and ignore
	var raw map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &raw)
	if err != nil {
		fmt.Printf("%s Failed to unmarshal debricked config: \"%s\"\n",
			color.YellowString("⚠️"),
			color.RedString(err.Error()),
		)

		return nil
	}

	// Accept any key for overrides, normalize to 'overrides'
	for k, v := range raw {
		lower := strings.ToLower(k)
		if lower == "overrides" || lower == "override" {
			raw["overrides"] = v
		}
	}

	// Marshal back to YAML and unmarshal into struct
	fixedYaml, err := yaml.Marshal(raw)
	if err != nil {
		fmt.Printf("%s Failed to re-marshal config: \"%s\"\n",
			color.YellowString("⚠️"),
			color.RedString(err.Error()),
		)

		return nil
	}

	err = yaml.Unmarshal(fixedYaml, &yamlConfig)
	if err != nil {
		fmt.Printf("%s Failed to unmarshal debricked config: \"%s\"\n",
			color.YellowString("⚠️"),
			color.RedString(err.Error()),
		)

		return nil
	}

	ignore := extractIgnore(raw)
	overrides := convertOverrides(yamlConfig.Overrides)

	return &DebrickedConfig{
		Overrides: overrides,
		Ignore:    ignore,
	}
}
