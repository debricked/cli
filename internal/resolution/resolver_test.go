package resolution

import (
	"errors"
	"fmt"
	"testing"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/file/testdata"
	resolutionFile "github.com/debricked/cli/internal/resolution/file"
	fileTestdata "github.com/debricked/cli/internal/resolution/file/testdata"
	"github.com/debricked/cli/internal/resolution/job"
	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"

	"github.com/debricked/cli/internal/resolution/strategy"
	strategyTestdata "github.com/debricked/cli/internal/resolution/strategy/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	workers   = 10
	goModFile = "go.mod"
)

func TestNewResolver(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)
	assert.NotNil(t, r)
}

func TestResolve(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	res, err := r.Resolve([]string{"../../go.mod"}, nil)
	assert.NotEmpty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestResolveInvokeError(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryErrorMock(),
		NewScheduler(workers),
	)

	_, err := r.Resolve([]string{"../../go.mod"}, nil)
	assert.NotNil(t, err)
}

func TestResolveStrategyError(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		fileTestdata.NewBatchFactoryMock(),
		strategy.NewStrategyFactory(),
		NewScheduler(workers),
	)

	res, err := r.Resolve([]string{"../../go.mod"}, nil)
	assert.Empty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestResolveScheduleError(t *testing.T) {
	errAssertion := errors.New("error")
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{Err: errAssertion},
	)

	res, err := r.Resolve([]string{"../../go.mod"}, nil)
	assert.NotEmpty(t, res.Jobs())
	assert.ErrorIs(t, err, errAssertion)
}

func TestResolveDirWithoutManifestFiles(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	res, err := r.Resolve([]string{"."}, nil)
	assert.Empty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestResolveInvalidDir(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	_, err := r.Resolve([]string{"invalid-dir"}, nil)
	assert.Error(t, err)
}

func TestResolveGetGroupsErr(t *testing.T) {
	f := testdata.NewFinderMock()
	testErr := errors.New("test")
	f.SetGetGroupsReturnMock(file.Groups{}, testErr)

	r := NewResolver(
		f,
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	_, err := r.Resolve([]string{"."}, nil)
	assert.ErrorIs(t, testErr, err)
}

func TestResolveDirWithManifestFiles(t *testing.T) {
	cases := []string{
		"",
		".",
		"./",
		"testdata",
		"./testdata/../testdata",
		"./strategy/testdata/",
		"strategy/testdata",
	}
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{ManifestFile: goModFile})
	f.SetGetGroupsReturnMock(groups, nil)

	r := NewResolver(
		f,
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	for _, dir := range cases {
		t.Run(fmt.Sprintf("Case: %s", dir), func(t *testing.T) {
			res, err := r.Resolve([]string{dir}, nil)
			assert.Len(t, res.Jobs(), 1)
			j := res.Jobs()[0]
			assert.False(t, j.Errors().HasError())
			assert.Equal(t, goModFile, j.GetFile())
			assert.NoError(t, err)
		})
	}
}

func TestResolveDirWithExclusions(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{ManifestFile: goModFile})
	f.SetGetGroupsReturnMock(groups, nil)

	r := NewResolver(
		f,
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	res, err := r.Resolve([]string{"."}, []string{"dir"})

	assert.Len(t, res.Jobs(), 1)
	j := res.Jobs()[0]
	assert.False(t, j.Errors().HasError())
	assert.Equal(t, goModFile, j.GetFile())
	assert.NoError(t, err)
}

func TestResolveHasResolutionErrs(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{ManifestFile: goModFile})
	f.SetGetGroupsReturnMock(groups, nil)

	jobErr := job.NewBaseJobError("job-error")
	jobWithErr := jobTestdata.NewJobMock(goModFile)
	jobWithErr.Errors().Warning(jobErr)
	schedulerMock := SchedulerMock{JobsMock: []job.IJob{jobWithErr}}

	r := NewResolver(
		f,
		resolutionFile.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		schedulerMock,
	)

	res, err := r.Resolve([]string{""}, []string{""})

	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 1)
	j := res.Jobs()[0]
	assert.Equal(t, goModFile, j.GetFile())
	assert.True(t, j.Errors().HasError())
	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.ErrorIs(t, jobErr, errs[0])
}
