package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

func MakePathFromManifestFile(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)
	if strings.EqualFold(string(os.PathSeparator), dir) {
		return fmt.Sprintf("%s%s", string(os.PathSeparator), fileName)
	}

	return fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)
}

func CloseFile(job job.IJob, fileWriter writer.IFileWriter, file *os.File) {
	err := fileWriter.Close(file)
	if err != nil {
		job.Errors().Critical(err)
	}
}
