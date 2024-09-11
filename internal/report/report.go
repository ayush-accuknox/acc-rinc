package report

import (
	"context"
	"io"
)

// Reporter defines an interface for reporting data. Implementations of this
// interface should write the report to the provided io.Writer and return any
// errors encountered during the process.
type Reporter interface {
	Report(ctx context.Context, to io.Writer) error
}
