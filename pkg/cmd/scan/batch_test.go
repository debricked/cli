package scan

import (
	"bytes"
	"debricked/pkg/client"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestConcludeWithoutAnyFiles(t *testing.T) {
	batch := newUploadBatch([]file.Group{}, nil)
	err := batch.conclude()
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if !strings.Contains(err.Error(), "failed to find dependency files") {
		t.Error("failed to asser error message")
	}
}

func TestUploadWithBadFiles(t *testing.T) {
	group := file.NewGroup("package.json", nil, []string{"yarn.lock"})
	fileGroups := []file.Group{
		*group,
	}
	metaObj, err := git.NewMetaObject("", "repository-name", "commit-name", "", "", "")
	if err != nil {
		t.Fatal("failed to create new MetaObject")
	}

	invalidToken := "invalid"
	debClient = client.NewDebClient(&invalidToken)
	batch := newUploadBatch(fileGroups, metaObj)
	output := captureOutput(batch.upload)
	outputAssertions := []string{
		"Failed to upload: package.json",
		"Unauthorized. Specify access token",
		"Failed to upload: yarn.lock",
	}
	for _, assertion := range outputAssertions {
		if !strings.Contains(output, assertion) {
			t.Error(fmt.Sprintf("failed to assert that output contained %s", assertion))
		}
	}

}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
