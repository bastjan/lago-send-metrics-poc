package lagoerr

import (
	"fmt"

	"github.com/getlago/lago-go-client"
)

type lagoErr struct {
	Err lago.Error
}

// Wrap wraps a lago.Error, which does not implement .Error(), into a standard error.
// Supports Unwrap().
func Wrap(err *lago.Error) error {
	if err == nil {
		return nil
	}
	return &lagoErr{Err: *err}
}

func (e *lagoErr) Error() string {
	if e == nil {
		return "<nil>"
	}
	wrappedMessage := "lago error"
	if e.Err.Err != nil {
		wrappedMessage = e.Err.Err.Error()
	}
	return fmt.Sprintf("%s: [%d] %s %+v", wrappedMessage, e.Err.HTTPStatusCode, e.Err.Msg, e.Err.ErrorDetail)
}

func (e *lagoErr) Unwrap() error {
	return e.Err.Err
}
