package upload

import (
	"debricked/pkg/client"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"errors"
)

type Options interface{}

type DebrickedOptions struct {
	FileGroups       file.Groups
	GitMetaObject    git.MetaObject
	IntegrationsName string
}

type IUploader interface {
	Upload(o Options) (*UploadResult, error)
}

type Uploader struct {
	client *client.IDebClient
}

func NewUploader(c *client.IDebClient) (*Uploader, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	return &Uploader{c}, nil
}

func (uploader *Uploader) Upload(o Options) (*UploadResult, error) {
	dOptions := o.(DebrickedOptions)
	batch := newUploadBatch(uploader.client, dOptions.FileGroups, &dOptions.GitMetaObject, dOptions.IntegrationsName)
	batch.upload()
	err := batch.conclude()
	if err != nil {
		return nil, err
	}

	result, err := batch.wait()
	if err != nil {
		return nil, err
	}

	return result, nil
}
