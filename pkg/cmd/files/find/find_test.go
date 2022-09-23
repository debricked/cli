package find

import (
	"debricked/pkg/file"
	"debricked/pkg/file/testdata"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewFindCmd(t *testing.T) {
	var f file.IFinder
	cmd := NewFindCmd(f)

	commands := cmd.Commands()
	nbrOfCommands := 0
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		if flag == nil {
			t.Error(fmt.Sprintf("failed to assert that %s flag was set", name))
		}
		if flag.Shorthand != shorthand {
			t.Error(fmt.Sprintf("failed to assert that %s flag shorthand %s was set correctly", name, shorthand))
		}
	}
}

func TestRunE(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(f)
	err := runE(nil, []string{"."})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestRunEError(t *testing.T) {
	f := testdata.NewFinderMock()
	errorAssertion := errors.New("finder-error")
	f.SetGetGroupsReturnMock(file.Groups{}, errorAssertion)
	runE := RunE(f)
	err := runE(nil, []string{"."})
	if err != errorAssertion {
		t.Fatal("failed to assert that error occured")
	}
}

func TestValidateArgs(t *testing.T) {
	err := validateArgs(nil, []string{"."})
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
}

func TestValidateArgsInvalidArgs(t *testing.T) {
	err := validateArgs(nil, []string{})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "requires path") {
		t.Error("failed to assert error message")
	}

	err = validateArgs(nil, []string{"invalid-path"})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "invalid path specified") {
		t.Error("failed to assert error message")
	}
}

//func (mock *debClientMock) Get(_ string, _ string) (*http.Response, error) {
//	var statusCode int
//	var body io.ReadCloser = nil
//	if clientMockAuthorized {
//		statusCode = http.StatusOK
//		formatsBytes, _ := json.Marshal(formatsMock)
//		body = ioutil.NewReadCloser(strings.NewReader(string(formatsBytes)), nil)
//	} else {
//		statusCode = http.StatusForbidden
//	}
//	res := http.Response{
//		Status:           "",
//		StatusCode:       statusCode,
//		Proto:            "",
//		ProtoMajor:       0,
//		ProtoMinor:       0,
//		Header:           nil,
//		Body:             body,
//		ContentLength:    0,
//		TransferEncoding: nil,
//		Close:            false,
//		Uncompressed:     false,
//		Trailer:          nil,
//		Request:          nil,
//		TLS:              nil,
//	}
//
//	return &res, nil
//}
