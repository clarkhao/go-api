package msgtype_test

import (
	"testing"

	msgtype "github.com/clarkhao/ws/utils/msgtype"
)

func TestGetMsg(t *testing.T) {
	cases := []struct {
		inputMsgType int
		inputMsg     []byte
		expected     string
	}{
		{
			inputMsgType: 1,
			inputMsg:     []byte("Hello"),
			expected:     "Hello",
		},
		{
			inputMsgType: 2,
			inputMsg:     []byte{72, 101, 108, 108, 111},
			expected:     "Hello",
		},
	}
	for _, c := range cases {
		output := msgtype.GetMsg(c.inputMsgType, c.inputMsg)
		if output != c.expected {
			t.Errorf("GetMsg(%v, %v)=%v, expected %v", c.inputMsgType, c.inputMsg, output, c.expected)
		}
	}
}
