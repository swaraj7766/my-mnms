package mnms

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/qeof/q"
)

type WebSocketMessage struct {
	Kind    string `json:"kind"`
	Level   int    `json:"level"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketStartWriteMessage() {

	for message := range QC.WebSocketMessageBroadcast {
		q.Q("got msg from WebSocketMessageBroadcast", len(QC.WebSocketClient))
		for client := range QC.WebSocketClient {
			if err := client.WriteJSON(message); err != nil {
				q.Q("error: websocket", err)
				continue
			}
			q.Q("websocket write to client", message)
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		q.Q(err)
	}
	QC.WebSocketClient[ws] = true
	q.Q("Client Connected")
	defer func() {
		delete(QC.WebSocketClient, ws)
		ws.Close()
		q.Q("Closed!")
	}()
	webSocketReader(ws)
}

func webSocketReader(conn *websocket.Conn) {
	for {
		var message WebSocketMessage
		err := conn.ReadJSON(&message)
		if !errors.Is(err, nil) {
			q.Q("error occurred: ", err)
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("error: ReadJSON %v", err)))
			if err != nil {
				q.Q("error while sending error message to ws client", err)
			}
			delete(QC.WebSocketClient, conn)
			break
		}
		q.Q(message)
	}
}
