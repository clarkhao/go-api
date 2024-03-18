package error

import (
	"errors"
	"fmt"
)

type RequestError struct {
	Code  int
	Msg   string
	Cause error
}

func (rs RequestError) Error() string {
	return fmt.Sprintf("StatusCode:%d, Message: %s, Cause: %v", rs.Code, rs.Msg, rs.Cause)
}

func (rs RequestError) Is(err error) bool {
	return errors.Is(rs.Cause, err) || rs.Cause.Error() == err.Error()
}
