package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"github.com/debricked/cli/pkg/tui"
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

var (
	NoFilesErr = errors.New("failed to find dependency files")
)

type uploadBatch struct {
	client          *client.IDebClient
	fileGroups      file.Groups
	gitMetaObject   *git.MetaObject
	integrationName string
	ciUploadId      int
}

func newUploadBatch(client *client.IDebClient, fileGroups file.Groups, gitMetaObject *git.MetaObject, integrationName string) *uploadBatch {
	return &uploadBatch{client: client, fileGroups: fileGroups, gitMetaObject: gitMetaObject, integrationName: integrationName, ciUploadId: 0}
}

// upload concurrently posts all file groups to Debricked
func (uploadBatch *uploadBatch) upload() {
	uploadWorker := func(fileQueue <-chan string, fileResults chan<- int) {
		const ok = 0
		const fail = 1
		for f := range fileQueue {
			err := uploadBatch.uploadFile(f)
			if err != nil {
				log.Println("Failed to upload:", f)
				if err != nil {
					log.Println(err.Error())
					fileResults <- fail
				}
			} else {
				fmt.Println("Successfully uploaded: ", f)
				fileResults <- ok
			}
		}
	}

	files := uploadBatch.fileGroups.GetFiles()
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
}

// uploadFile Reads file content from filepath and uploads it to Debricked. Returns HTTP status code or 0 if other error occur
func (uploadBatch *uploadBatch) uploadFile(filePath string) error {
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

	_ = writer.WriteField("fileRelativePath", filepath.Dir(filePath))
	_ = writer.WriteField("repositoryName", uploadBatch.gitMetaObject.RepositoryName)
	_ = writer.WriteField("commitName", uploadBatch.gitMetaObject.CommitName)
	_ = writer.WriteField("repositoryUrl", uploadBatch.gitMetaObject.RepositoryUrl)
	_ = writer.WriteField("branchName", uploadBatch.gitMetaObject.BranchName)
	_ = writer.WriteField("defaultBranchName", uploadBatch.gitMetaObject.DefaultBranchName)
	if uploadBatch.initialized() {
		_ = writer.WriteField("ciUploadId", strconv.Itoa(uploadBatch.ciUploadId))
	}
	response, err := (*uploadBatch.client).Post(
		"/api/1.0/open/uploads/dependencies/files",
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return err
	}

	mutex := sync.Mutex{}
	mutex.Lock()
	if !uploadBatch.initialized() {
		data, _ := io.ReadAll(response.Body)
		defer response.Body.Close()

		uFile := uploadedFile{}
		_ = json.Unmarshal(data, &uFile)
		uploadBatch.ciUploadId = uFile.CiUploadId
	}
	mutex.Unlock()

	return nil
}

// conclude send the conclusion request to Debricked
func (uploadBatch *uploadBatch) conclude() error {
	if uploadBatch.ciUploadId == 0 {
		return NoFilesErr
	}
	body, err := json.Marshal(uploadConclusion{
		CiUploadId:      strconv.Itoa(uploadBatch.ciUploadId),
		RepositoryName:  uploadBatch.gitMetaObject.RepositoryName,
		IntegrationName: uploadBatch.integrationName,
		CommitName:      uploadBatch.gitMetaObject.CommitName,
		Author:          uploadBatch.gitMetaObject.Author,
	})

	if err != nil {
		return err
	}
	response, err := (*uploadBatch.client).Post(
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
		if res.StatusCode == http.StatusCreated {
			err := bar.Finish()
			if err != nil {
				return nil, err
			}
			return nil, errors.New("progress polling terminated due to long queue times")
		}
		status, err := newUploadStatus(res)
		if err != nil {
			return nil, err
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
