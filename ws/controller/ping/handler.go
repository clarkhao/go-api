package ping

import (
	"fmt"
	"log"
	"time"

	msgtype "github.com/clarkhao/ws/utils/msgtype"
	ws "github.com/clarkhao/ws/utils/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WSConnHandler(c *gin.Context) {
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		//extract the data received
		if messageType == websocket.CloseMessage {
			conn.Close()
			return
		}

		if msgtype.IsJson(message) {
			data := msgtype.GetMsgObj[msgtype.THello](message)
			fmt.Println("Received:", data.Hello)
		} else {
			data := msgtype.GetMsg(messageType, message)
			fmt.Println("Received:", data, messageType)
		}
		conn.WriteMessage(messageType, message)
		time.Sleep(time.Second)
	}
}
