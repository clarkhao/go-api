package error_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	re "restapi/utils/error"
)

func TestRequestError(t *testing.T) {
	cases := []struct {
		name            string
		code            int
		msg             string
		cause           error
		expectedWrapped bool
	}{
		{
			name:            "client error",
			code:            400,
			msg:             "from some Get Handler",
			cause:           fmt.Errorf("%w", http.ErrAbortHandler),
			expectedWrapped: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := re.RequestError{
				Code:  tc.code,
				Msg:   tc.msg,
				Cause: tc.cause,
			}
			ok := errors.Is(err, http.ErrAbortHandler)
			if !ok {
				t.Errorf("failed to wrapped error")
			}
		})
	}
}
