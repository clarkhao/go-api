package ping

import (
	"errors"

	re "restapi/utils/error"
)

func SomeServiceErr() error {
	return re.RequestError{
		Code:  401,
		Msg:   "got authentication err From ping/SomeServiceErr",
		Cause: errors.New("Authentication Error"),
	}
}
