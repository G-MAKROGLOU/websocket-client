package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

var wg sync.WaitGroup

// New creates a new SocketClient
func New(origin string, server string, events SocketClientEvents) *SocketClient {
	return &SocketClient{
		Origin: origin,
		Server: server,
		events: events,
	}
}

// Connect connects a client to a socket server
func (sc *SocketClient) Connect() error {
	id := uuid.NewString()

	config, _ := websocket.NewConfig(sc.Server, sc.Origin)

	config.Header = http.Header{
		"Cookie": {"session_id=" + id},
	}

	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}

	sc.Conn = ws
	sc.ID = id

	msg := fmt.Sprintf("[client] connected with session id: %s", id)
	slog.Info(msg)

	return nil
}

// Disconnect disconnects from the server
func (sc *SocketClient) Disconnect() {
	if err := sc.sendDisconnect(); err != nil {
		msg := fmt.Sprintf("disconnect: %v", err)
		slog.Error(msg, "error", err.Error())
	}
}

func (sc *SocketClient) sendDisconnect() error {
	data := map[string]interface{}{
		"GmWsType": "disconnect",
	}
	return websocket.JSON.Send(sc.Conn, data)
}

// ReceiveJSON handles incoming messages. To be started in a goroutine
func (sc *SocketClient) ReceiveJSON() {
	for {
		var data map[string]interface{}
		err := websocket.JSON.Receive(sc.Conn, &data)

		// no errors, continue receiving
		if err == nil {
			sc.events.OnReceive(data)
			continue
		}

		// io.EOF, connection closed
		if err == io.EOF {
			break
		}

		// other errors, let the user handle them
		if err != io.EOF {
			sc.events.OnReceiveError(err)
			continue
		}
	}
}

// ReceiveText handles text incoming messages that you can deserialize to whatver type you want
// in the onReceive event. Suitable for connections to servers that are not created with the github.com/G-MAKROGLOU/websocket-server
// package.
func (sc *SocketClient) ReceiveText() {
	for {
		var data map[string]interface{}
		var buff = make([]byte, 4096)

		size, err := sc.Conn.Read(buff)

		// io.EOF, connection closed
		if err == io.EOF {
			break
		}

		// other errors, let the user handle them
		if err != nil && err != io.EOF {
			sc.events.OnReceiveError(err)
			continue
		}

		uErr := json.Unmarshal(buff[:size], &data)

		// unmarshall error, let the user know, possible buffer resizing error
		if uErr != nil {
			sc.events.OnReceiveError(err)
			continue
		}

		// all good
		sc.events.OnReceive(data)
	}
}

// Join adds a socket to a room (1-N message exchange)
func (sc *SocketClient) Join(roomName string) {
	data := map[string]interface{}{
		"GmWsType": "join",
		"GmWsRoom": roomName,
	}
	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnJoinError(roomName, err)
		return
	}
}

// Leave leaves a room
func (sc *SocketClient) Leave(roomName string) {
	data := map[string]interface{}{
		"GmWsType": "leave",
		"GmWsRoom": roomName,
	}
	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnLeaveError(roomName, err)
		return
	}
}

// SendJSON sends a broadcast message to all connected sockets on the server. Can be used with any server that supports JSON, but it will add an extra property
// in case it is used with github.com/G-MAKROGLOU/websocket-server that supports rooms.
func (sc *SocketClient) SendJSON(data map[string]interface{}) {
	data["GmWsType"] = "broadcast"

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
}

// SendJSONTo sends a unicast/multicast message to all sockets in a rooml. Can be used with any server that supports JSON, but it will add an extra property
// in case it is used with github.com/G-MAKROGLOU/websocket-server that supports rooms.
func (sc *SocketClient) SendJSONTo(roomName string, data map[string]interface{}) {
	data["GmWsType"] = "multicast"
	data["GmWsRoom"] = roomName

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
}

// SendText sends text over the wire. Can be used with any socket server that comminicates with text. Not supported by github.com/G-MAKROGLOU/websocket-server
func (sc *SocketClient) SendText(msg interface{}) {

	b, _ := json.Marshal(msg)

	_, err := sc.Conn.Write(b)

	if err != nil {
		sc.events.OnSendError(err)
		return
	}
	var data map[string]interface{}
	json.Unmarshal(b, &data)
}
