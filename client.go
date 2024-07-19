package client

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

var wg sync.WaitGroup

// SocketClient represents a socket client connection 
type SocketClient struct {
	ID string
	Origin string
	Server string
	Conn *websocket.Conn
	events SocketClientEvents
}

// SocketClientEvents holds all the possible events that are supported
type SocketClientEvents interface {
	OnConnect(ws *websocket.Conn, sessID string)
	OnConnectError(err error)
	OnDisconnect()
	OnDisconnectError(err error)
	OnReceive(data map[string]interface{})
	OnReceiveError(err error)
	OnJoin(roomName string)
	OnJoinError(roomName string, err error)
	OnLeave(roomName string)
	OnLeaveError(roomName string, err error)
	OnSend(data map[string]interface{})
	OnSendError(err error)
}

// NewClient creates a new SocketClient
func NewClient(origin string, server string, events SocketClientEvents) *SocketClient {
	return &SocketClient{
		Origin: origin,
		Server: server,
		events: events,
	}
}

//Connect connects a client to a socket server
func (sc *SocketClient) Connect() error {	
	id := uuid.NewString()
	
	config, _ := websocket.NewConfig(sc.Server, sc.Origin)

	config.Header = http.Header{
        "Cookie": {"session_id=" + id},
    }
	
	ws, err := websocket.DialConfig(config)
	if err != nil {
		sc.events.OnConnectError(err)
		return err
	}

	sc.Conn = ws
	sc.ID = id
	sc.events.OnConnect(ws, id)

	return nil
}

// Disconnect disconnects from the server
func (sc *SocketClient) Disconnect() {
	data := map[string]interface{}{
		"Type": "gm_ws_disconnect",
	}
	sc.SendJSON(data)
	if err := sc.Conn.Close(); err != nil {
		sc.events.OnDisconnectError(err)
		return
	}
	sc.events.OnDisconnect()
}

// ReceiveJSON handles incoming messages. To be started in a goroutine
func (sc *SocketClient) ReceiveJSON() {
	for {
		var data map[string]interface{}
		if err := websocket.JSON.Receive(sc.Conn, &data); err != nil {
			sc.events.OnReceiveError(err)
			break
		}
		sc.events.OnReceive(data)
	}
}

// ReceiveText handles text incoming messages that you can deserialize to whatver type you want
// in the OnReceive event. Suitable for connections to servers that are not created with the github.com/G-MAKROGLOU/websocket-server 
// package.
func (sc *SocketClient) ReceiveText() {
	for {
		var data map[string]interface{}
		var buff = make([]byte, 4096)

		size, err := sc.Conn.Read(buff)
		if err != nil {
			sc.events.OnReceiveError(err)
			continue
		}

		umErr := json.Unmarshal(buff[:size], &data)
		if umErr != nil {
			sc.events.OnReceiveError(err)
		}

		sc.events.OnReceive(data)
	}
}

// Join adds a socket to a room (1-N message exchange)
func (sc *SocketClient) Join(roomName string) {
	data := map[string]interface{}{
		"Gm_Ws_Type": "gm_ws_join",
		"Gm_Ws_Room": roomName,
	}
	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnJoinError(roomName, err)
		return
	}
	sc.events.OnJoin(roomName)
}

// Leave leaves a room
func (sc *SocketClient) Leave(roomName string) {
	data := map[string]interface{}{
		"Gm_Ws_Type": "gm_ws_leave",
		"Gm_Ws_Room": roomName,
	}
	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnLeaveError(roomName, err)
		return
	}
	sc.events.OnLeave(roomName)
}

// SendJSON sends a broadcast message to all connected sockets on the server. Can be used with any server that supports JSON, but it will add an extra property
// in case it is used with github.com/G-MAKROGLOU/websocket-server that supports rooms.
func (sc *SocketClient) SendJSON(data map[string]interface{}) {
	data["Gm_Ws_Type"] = "gm_ws_broadcast"

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
	sc.events.OnSend(data)
}

// SendJSONTo sends a unicast/multicast message to all sockets in a rooml. Can be used with any server that supports JSON, but it will add an extra property
// in case it is used with github.com/G-MAKROGLOU/websocket-server that supports rooms.
func (sc *SocketClient) SendJSONTo(roomName string, data map[string]interface{}) {
	data["Gm_Ws_Type"] = "gm_ws_multicast"
	data["Gm_Ws_Room"] = roomName

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
	sc.events.OnSend(data)
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
	sc.events.OnSend(data)
}






