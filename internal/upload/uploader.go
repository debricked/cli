package upload

import (
	"errors"

	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/git"
)

type IOptions interface{}

type DebrickedOptions struct {
	FileGroups             file.Groups
	GitMetaObject          git.MetaObject
	IntegrationsName       string
	CallGraphUploadTimeout int
	VersionConsolidation   int
}

type IUploader interface {
	Upload(o IOptions) (*UploadResult, error)
}

type Uploader struct {
	client *client.IDebClient
}

func NewUploader(c client.IDebClient) (*Uploader, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}

	return &Uploader{&c}, nil
}

func (uploader *Uploader) Upload(o IOptions) (*UploadResult, error) {
	dOptions := o.(DebrickedOptions)
	batch := newUploadBatch(uploader.client, dOptions.FileGroups, &dOptions.GitMetaObject, dOptions.IntegrationsName, dOptions.CallGraphUploadTimeout, dOptions.VersionConsolidation)

	err := batch.upload()
	if err != nil {
		return nil, err
	}

	err = batch.initAnalysis()
	if err != nil {
		return nil, err
	}

	result, err := batch.wait()
	if err != nil {
		// the command should not fail because some file can't be scanned
		if err == PollingTerminatedErr {
			return result, nil
		}

		return nil, err
	}

	return result, nil
}
