package docker

import (
	"fmt"
	"strings"
)

const tagDelimiter = ":"

type Image struct {
	Repository string
	Tag        string
	Digest     string
}

func (i Image) Name() string {
	return fmt.Sprintf("%s:%s", i.Repository, i.Tag)
}

func MakeImage(image string) Image {
	var img Image
	var repository string
	var tag string
	hasTag := strings.Contains(image, tagDelimiter)

	if hasTag {
		parts := strings.Split(image, tagDelimiter)
		repository = parts[0]
		tag = parts[1]
	} else {
		repository = image
		tag = "latest"
	}

	img.Repository = repository
	img.Tag = tag

	return img
}
