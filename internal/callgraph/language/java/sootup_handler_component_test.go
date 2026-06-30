package java

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSootUpHandlerComponent_GetSootUpWrapperInvalidVersionString(t *testing.T) {
	h := SootUpHandler{}
	_, err := h.GetSootUpWrapper("not-a-number", ioTestData.FileSystemMock{}, ioTestData.ArchiveMock{})
	assert.Error(t, err)
}

func TestSootUpHandlerComponent_GetSootUpWrapperUnsupportedJavaVersion(t *testing.T) {
	h := SootUpHandler{}
	_, err := h.GetSootUpWrapper("8", ioTestData.FileSystemMock{}, ioTestData.ArchiveMock{})
	assert.Error(t, err)
}

func TestSootUpHandlerComponent_GetSootUpWrapperMkdirError(t *testing.T) {
	h := SootUpHandler{}
	errString := "mkdir error"
	fsMock := ioTestData.FileSystemMock{
		IsNotExistBool: true,
		MkdirError:     fmt.Errorf("%s", errString), //nolint
	}
	arcMock := ioTestData.ArchiveMock{}
	_, err := h.GetSootUpWrapper("21", fsMock, arcMock)
	assert.EqualError(t, err, errString)
}

func TestSootUpHandlerComponent_GetSootUpWrapperVersion21UsesEmbeddedJar(t *testing.T) {
	h := SootUpHandler{}
	fsMock := ioTestData.FileSystemMock{
		StatError:      fmt.Errorf("missing"),
		IsNotExistBool: true,
	}

	p, err := h.GetSootUpWrapper("21", fsMock, ioTestData.ArchiveMock{})
	assert.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.True(t, strings.HasSuffix(filepath.ToSlash(p), "/.debricked/SootUpWrapper.jar"))
}

func TestSootUpHandlerComponent_GetSootUpWrapperReturnsExistingPath(t *testing.T) {
	h := SootUpHandler{}
	fsMock := ioTestData.FileSystemMock{}

	p, err := h.GetSootUpWrapper("17", fsMock, ioTestData.ArchiveMock{})
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(filepath.ToSlash(p), "/.debricked/SootUpWrapper.jar"))
}
