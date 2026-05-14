package java

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestSootHandlerComponent_GetSootWrapperInvalidVersionString(t *testing.T) {
	h := SootHandler{}
	_, err := h.GetSootWrapper("not-a-number", ioTestData.FileSystemMock{}, ioTestData.ArchiveMock{})
	assert.EqualError(t, err, "could not convert version to int")
}

func TestSootHandlerComponent_GetSootWrapperUnsupportedJavaVersion(t *testing.T) {
	h := SootHandler{}
	_, err := h.GetSootWrapper("8", ioTestData.FileSystemMock{}, ioTestData.ArchiveMock{})
	assert.EqualError(t, err, "lowest supported version for running callgraph generation is 11")
}

func TestSootHandlerComponent_GetSootWrapperMkdirError(t *testing.T) {
	h := SootHandler{}
	fsMock := ioTestData.FileSystemMock{
		StatError:      fmt.Errorf("missing"),
		IsNotExistBool: true,
		MkdirError:     fmt.Errorf("mkdir failed"),
	}

	_, err := h.GetSootWrapper("11", fsMock, ioTestData.ArchiveMock{})
	assert.EqualError(t, err, "mkdir failed")
}

func TestSootHandlerComponent_GetSootWrapperVersion21UsesEmbeddedJar(t *testing.T) {
	h := SootHandler{}
	fsMock := ioTestData.FileSystemMock{
		StatError:      fmt.Errorf("missing"),
		IsNotExistBool: true,
	}

	p, err := h.GetSootWrapper("21", fsMock, ioTestData.ArchiveMock{})
	assert.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.True(t, strings.HasSuffix(filepath.ToSlash(p), "/.debricked/SootWrapper.jar"))
}

func TestSootHandlerComponent_GetSootWrapperReturnsExistingPath(t *testing.T) {
	h := SootHandler{}
	fsMock := ioTestData.FileSystemMock{}

	p, err := h.GetSootWrapper("17", fsMock, ioTestData.ArchiveMock{})
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(filepath.ToSlash(p), "/.debricked/soot-wrapper.jar"))
}

