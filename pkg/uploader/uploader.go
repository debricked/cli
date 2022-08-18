package uploader

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

type Uploader interface {
	Upload(o Options) (*UploadResult, error)
}

type debrickedUploader struct {
	client *client.Client
}

func NewDebrickedUploader(c *client.Client) (*debrickedUploader, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	return &debrickedUploader{c}, nil
}

func (uploader *debrickedUploader) Upload(o Options) (*UploadResult, error) {
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
