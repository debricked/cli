package testdata

import "github.com/debricked/cli/internal/report"

type ReporterMock struct {
	err error
}

func NewReporterMock() *ReporterMock {
	return &ReporterMock{nil}
}

func (r *ReporterMock) Order(_ report.IOrderArgs) error {
	return r.err
}

func (r *ReporterMock) SetError(e error) {
	r.err = e
}
