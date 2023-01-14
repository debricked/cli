package scan

import (
	"errors"
)

var (
	BadOptsErr = errors.New("failed to type case IOptions")
)

type IScanner interface {
	Scan(o IOptions) error
}
