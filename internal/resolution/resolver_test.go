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
	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{"../../go.mod"}, options)
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
	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	_, err := r.Resolve([]string{"../../go.mod"}, options)
	assert.NotNil(t, err)
}

func TestResolveStrategyError(t *testing.T) {
	r := NewResolver(
		&testdata.FinderMock{},
		fileTestdata.NewBatchFactoryMock(),
		strategy.NewStrategyFactory(),
		NewScheduler(workers),
	)

	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{"../../go.mod"}, options)
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

	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{"../../go.mod"}, options)
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

	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{"."}, options)
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

	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	_, err := r.Resolve([]string{"invalid-dir"}, options)
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

	options := DebrickedOptions{
		Exclusions: nil,
		Verbose:    true,
		Regenerate: 0,
	}
	_, err := r.Resolve([]string{"."}, options)
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

	for i, dir := range cases {
		options := DebrickedOptions{
			Exclusions: nil,
			Verbose:    true,
			Regenerate: i % 3, // To test the different regenerate values
		}
		t.Run(fmt.Sprintf("Case: %s", dir), func(t *testing.T) {
			res, err := r.Resolve([]string{dir}, options)
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

	options := DebrickedOptions{
		Exclusions: []string{"dir"},
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{"."}, options)

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

	options := DebrickedOptions{
		Exclusions: []string{""},
		Verbose:    true,
		Regenerate: 0,
	}
	res, err := r.Resolve([]string{""}, options)

	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 1)
	j := res.Jobs()[0]
	assert.Equal(t, goModFile, j.GetFile())
	assert.True(t, j.Errors().HasError())
	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.ErrorIs(t, jobErr, errs[0])
}

func TestGetExitCode(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{ManifestFile: goModFile})
	f.SetGetGroupsReturnMock(groups, nil)

	jobErr := job.NewBaseJobError("job-error")
	jobWithErr := jobTestdata.NewJobMock(goModFile)
	jobWithErr.Errors().Warning(jobErr)

	jobNoErr := jobTestdata.NewJobMock(goModFile)

	cases := []struct {
		strictness       StrictnessLevel
		expected         int
		nbFailingJobs    int
		nbSuccessfulJobs int
	}{
		{NoFail, 0, 1, 1},
		{NoFail, 0, 1, 0},
		{FailIfAllFail, 0, 1, 1},
		{FailIfAllFail, 1, 2, 0},
		{FailIfAllFail, 0, 0, 2},
		{FailIfAnyFail, 1, 1, 1},
		{FailIfAnyFail, 1, 2, 0},
		{FailIfAnyFail, 0, 0, 2},
		{FailOrWarn, 3, 1, 1},
		{FailOrWarn, 1, 2, 0},
		{FailOrWarn, 0, 0, 2},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Strictness: %d", c.strictness), func(t *testing.T) {
			var jobs []job.IJob
			for i := 0; i < c.nbFailingJobs; i++ {
				jobs = append(jobs, jobWithErr)
			}
			for i := 0; i < c.nbSuccessfulJobs; i++ {
				jobs = append(jobs, jobNoErr)
			}
			schedulerMock := SchedulerMock{JobsMock: jobs}
			r := NewResolver(
				f,
				resolutionFile.NewBatchFactory(),
				strategyTestdata.NewStrategyFactoryMock(),
				schedulerMock,
			)
			options := DebrickedOptions{
				Exclusions:           []string{""},
				Verbose:              true,
				Regenerate:           0,
				Resolutionstrictness: c.strictness,
			}
			resolution, err := r.Resolve([]string{""}, options)
			if c.expected > 0 {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			exitCode, err := r.GetExitCode(resolution, options)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, exitCode)
		})
	}
}

func TestGetStrictnessLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		want    StrictnessLevel
		wantErr bool
	}{
		{
			name:    "Test NoFail",
			level:   0,
			want:    NoFail,
			wantErr: false,
		},
		{
			name:    "Test FailIfAllFail",
			level:   1,
			want:    FailIfAllFail,
			wantErr: false,
		},
		{
			name:    "Test FailIfAnyFail",
			level:   2,
			want:    FailIfAnyFail,
			wantErr: false,
		},
		{
			name:    "Test FailOrWarn",
			level:   3,
			want:    FailOrWarn,
			wantErr: false,
		},
		{
			name:    "Test Invalid",
			level:   151,
			want:    NoFail,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStrictnessLevel(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStrictnessLevel() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetStrictnessLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
