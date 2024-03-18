package msgtype

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type THello struct {
	Hello string
}

func IsJson(msg []byte) bool {
	prefix := strings.HasPrefix(string(msg), "{")
	suffix := strings.HasSuffix(string(msg), "}")
	return prefix && suffix
}

func GetMsg(t int, msg []byte) string {
	switch t {
	case websocket.TextMessage:
		return string(msg)
	case websocket.BinaryMessage:
		base64Str := base64.StdEncoding.EncodeToString(msg)
		decoded, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			fmt.Println("Error decoding Base64:", err)
		}
		return string(decoded)
	default:
		return string(msg)
	}
}

func GetMsgObj[T interface{}](msg []byte) (value T) {
	json.Unmarshal(msg, &value)
	return
}

func SendMsg(t int, msg []byte) []byte {
	switch t {
	case websocket.TextMessage:
		return msg
	case websocket.BinaryMessage:
		return []byte(base64.StdEncoding.EncodeToString(msg))
	default:
		return msg
	}
}
