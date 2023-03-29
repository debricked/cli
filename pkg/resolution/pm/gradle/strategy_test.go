package gradle

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{"file"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)

	s = NewStrategy([]string{"file-1", "file-2"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy([]string{}, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	s := NewStrategy([]string{"file"}, nil)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	s := NewStrategy([]string{"file-1", "file-2"}, nil)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 2)
}

// mock for IGradleSetup
type mockGradleSetup struct {
	mock.Mock
}

// mock for Setup
func (m *mockGradleSetup) Setup(files []string, paths []string) (GradleSetup, error) {
	args := m.Called()

	return args.Get(0).(GradleSetup), args.Error(1)
}

func TestInvokeError(t *testing.T) {
	s := NewStrategy([]string{"file"}, []string{"path"})
	mocked := &mockGradleSetup{}
	mocked.On("Setup").Return(GradleSetup{}, GradleSetupWalkError{})

	s.GradleSetup = mocked
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeFoundProject(t *testing.T) {
	s := NewStrategy([]string{"file"}, []string{"file"})
	subprojectMap := make(map[string]string)
	dir, _ := os.Getwd()
	subprojectMap[dir] = ""
	mocked := &mockGradleSetup{}
	mocked.On("Setup").Return(GradleSetup{GradleProjects: []GradleProject{{dir: dir, gradlew: "gradlew"}}, groovyScriptPath: "", subProjectMap: subprojectMap}, nil)

	s.GradleSetup = mocked
	jobs, _ := s.Invoke()

	assert.Len(t, jobs, 1)
}
