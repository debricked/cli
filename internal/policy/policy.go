package policy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/client"
	internalIO "github.com/debricked/cli/internal/io"
)

// DefaultFileName is the policy file looked for when no path is given
const DefaultFileName = "debricked_policy.json"

const validateUri = "/api/1.0/open/debricked-policy/validate"

var ErrNoPolicyFile = errors.New("no policy file found")

type IValidator interface {
	// Validate checks a Debricked policy file against Debricked's policy schema.
	// An empty path falls back to DefaultFileName in the working directory.
	Validate(path string) (Result, error)
}

// Result is the outcome of validating one policy file
type Result struct {
	File   string            `json:"file"`
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

type Validator struct {
	DebClient  client.IDebClient
	FileSystem internalIO.IFileSystem
}

type validationRequest struct {
	Content string `json:"content"`
}

func (v Validator) Validate(path string) (Result, error) {
	filePath, err := v.resolvePath(path)
	if err != nil {
		return Result{}, err
	}

	content, err := v.FileSystem.ReadFile(filePath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to read %s. %s", filePath, err)
	}

	result := Result{File: filePath, Errors: []ValidationError{}}

	// Malformed JSON never makes it past the schema validation, so report it
	// without spending a request on it.
	if syntaxErr := jsonSyntaxError(content); syntaxErr != nil {
		result.Errors = append(result.Errors, *syntaxErr)

		return result, nil
	}

	validationErrors, err := v.validateContent(content)
	if err != nil {
		return Result{}, err
	}
	result.Errors = validationErrors
	result.Valid = len(validationErrors) == 0

	return result, nil
}

// resolvePath returns the policy file to validate. Directories are expanded
// with DefaultFileName, which also covers the empty, defaulted path.
func (v Validator) resolvePath(path string) (string, error) {
	if len(path) == 0 {
		path = DefaultFileName
	}

	fileInfo, err := v.statPolicyFile(path)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return path, nil
	}

	path = filepath.Join(path, DefaultFileName)
	if _, err = v.statPolicyFile(path); err != nil {
		return "", err
	}

	return path, nil
}

func (v Validator) statPolicyFile(path string) (os.FileInfo, error) {
	stat, err := v.FileSystem.Stat(path)
	if err != nil {
		if v.FileSystem.IsNotExist(err) {
			return nil, fmt.Errorf(
				"%w at %s. Run `debricked policy validate <path>` to validate a file in another location",
				ErrNoPolicyFile,
				path,
			)
		}

		return nil, fmt.Errorf("failed to read %s. %s", path, err)
	}

	return stat, nil
}

func (v Validator) validateContent(content []byte) ([]ValidationError, error) {
	body, err := json.Marshal(validationRequest{Content: string(content)})
	if err != nil {
		return nil, fmt.Errorf("failed to create validation request. %s", err)
	}

	res, err := v.DebClient.Post(validateUri, "application/json", bytes.NewBuffer(body), client.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation response. %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to validate policy file. Status code: %d %s",
			res.StatusCode,
			strings.TrimSpace(string(data)),
		)
	}

	var validationErrors []ValidationError
	if err = json.Unmarshal(data, &validationErrors); err != nil {
		return nil, fmt.Errorf("failed to parse validation response. %s", err)
	}

	if validationErrors == nil {
		validationErrors = []ValidationError{}
	}

	return validationErrors, nil
}
