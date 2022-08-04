package scan

import (
	"bytes"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type uploadBatch struct {
	fileGroups    file.Groups
	gitMetaObject *git.MetaObject
	ciUploadId    int
}

// upload concurrently posts all file groups to Debricked
func (uploadBatch *uploadBatch) upload() {
	var wg sync.WaitGroup

	uploadWorker := func(filePath string) {
		// Mark upload done
		if uploadBatch.initialized() {
			defer wg.Done()
		}
		err := uploadBatch.uploadFile(filePath)
		if err != nil {
			log.Println("Failed to upload:", filePath)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			fmt.Println("Successfully uploaded: ", filePath)
		}
	}

	for _, f := range uploadBatch.fileGroups.GetFiles() {
		if !uploadBatch.initialized() {
			uploadWorker(f)
		} else {
			// Increment WaitGroup Counter
			wg.Add(1)
			go uploadWorker(f)
		}
	}
	// Wait for goroutines to finish
	wg.Wait()
}

func newUploadBatch(fileGroups file.Groups, gitMetaObject *git.MetaObject) *uploadBatch {
	return &uploadBatch{fileGroups: fileGroups, gitMetaObject: gitMetaObject, ciUploadId: 0}
}

// uploadFile Reads file content from filepath and uploads it to Debricked. Returns HTTP status code or 0 if other error occur
func (uploadBatch *uploadBatch) uploadFile(filePath string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileData, _ := writer.CreateFormFile("fileData", filepath.Base(filePath))
	f, _ := os.Open(filePath)
	defer f.Close()
	_, _ = io.Copy(fileData, f)

	_ = writer.WriteField("fileRelativePath", filepath.Dir(filePath))
	_ = writer.WriteField("repositoryName", uploadBatch.gitMetaObject.RepositoryName)
	_ = writer.WriteField("commitName", uploadBatch.gitMetaObject.CommitName)
	_ = writer.WriteField("repositoryUrl", uploadBatch.gitMetaObject.RepositoryUrl)
	_ = writer.WriteField("branchName", uploadBatch.gitMetaObject.BranchName)
	_ = writer.WriteField("defaultBranchName", uploadBatch.gitMetaObject.DefaultBranchName)
	if uploadBatch.initialized() {
		_ = writer.WriteField("ciUploadId", strconv.Itoa(uploadBatch.ciUploadId))
	}
	response, err := debClient.Post(
		"/api/1.0/open/uploads/dependencies/files",
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return err
	}

	if !uploadBatch.initialized() {
		data, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		uploadedFile := uploadedFile{}
		_ = json.Unmarshal(data, &uploadedFile)
		uploadBatch.ciUploadId = uploadedFile.CiUploadId
	}

	return nil
}

// conclude send the conclusion request to Debricked
func (uploadBatch *uploadBatch) conclude() error {
	if uploadBatch.ciUploadId == 0 {
		return errors.New("failed to find dependency files")
	}
	body, err := json.Marshal(uploadConclusion{
		CiUploadId:      strconv.Itoa(uploadBatch.ciUploadId),
		RepositoryName:  uploadBatch.gitMetaObject.RepositoryName,
		IntegrationName: integrationName,
		CommitName:      uploadBatch.gitMetaObject.CommitName,
		Author:          uploadBatch.gitMetaObject.Author,
	})

	if err != nil {
		return err
	}
	response, err := debClient.Post(
		"/api/1.0/open/finishes/dependencies/files/uploads",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Failed to conclude upload due to status code %d", response.StatusCode))
	} else {
		fmt.Println("Successfully concluded upload")
	}

	return nil
}

func (uploadBatch *uploadBatch) initialized() bool {
	return uploadBatch.ciUploadId > 0
}

// wait track scan progress and return scanStatus upon completion
func (uploadBatch *uploadBatch) wait() (*scanStatus, error) {
	bar := newProgressBar()
	_ = bar.RenderBlank()
	// poll scan status until completion
	var resultStatus *scanStatus
	uri := fmt.Sprintf("/api/1.0/open/ci/upload/status?ciUploadId=%s", strconv.Itoa(uploadBatch.ciUploadId))
	for !bar.IsFinished() {
		res, err := debClient.Get(uri, "application/json")
		if err != nil {
			return nil, err
		}
		if res.StatusCode == http.StatusCreated {
			err := bar.Finish()
			if err != nil {
				return nil, err
			}
			return nil, errors.New("progress polling terminated due to long queue times")
		}
		status, err := newScanStatus(res)
		if err != nil {
			return nil, err
		}
		err = bar.Set(status.Progress)
		if err != nil {
			return nil, err
		}

		if bar.IsFinished() {
			resultStatus = status
		} else {
			time.Sleep(2000 * time.Millisecond)
		}
	}

	return resultStatus, nil
}

type uploadedFile struct {
	CiUploadId           int    `json:"ciUploadId"`
	UploadProgramsFileId int    `json:"uploadProgramsFileId"`
	TotalScans           int    `json:"totalScans"`
	RemainingScans       int    `json:"remainingScans"`
	Percentage           string `json:"percentage"`
	EstimateDaysLeft     int    `json:"estimateDaysLeft"`
}

type uploadConclusion struct {
	CiUploadId      string `json:"ciUploadId"`
	RepositoryName  string `json:"repositoryName"`
	IntegrationName string `json:"integrationName"`
	CommitName      string `json:"commitName"`
	Author          string `json:"author"`
}

func newProgressBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[blue]Scanning...[reset]"),
		progressbar.OptionOnCompletion(func() {
			color.NoColor = false
			color.Green("✔")
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)
}
