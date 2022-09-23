package testdata

import "debricked/pkg/report"

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
