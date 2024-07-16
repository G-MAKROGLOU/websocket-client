package client

import (
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

// NewSocketClient creates a new SocketClient
func NewSocketClient(origin string, server string, events SocketClientEvents) *SocketClient {
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
	sc.Send(data)
	if err := sc.Conn.Close(); err != nil {
		sc.events.OnDisconnectError(err)
		return
	}
	sc.events.OnDisconnect()
}

// Receive handles incoming messages. To be started in a goroutine
func (sc *SocketClient) Receive() {
	for {
		var data map[string]interface{}
		if err := websocket.JSON.Receive(sc.Conn, &data); err != nil {
			sc.events.OnReceiveError(err)
			break
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

// Send sends a broadcast message to all connected sockets on the server
func (sc *SocketClient) Send(data map[string]interface{}) {
	data["Gm_Ws_Type"] = "gm_ws_broadcast"

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
	sc.events.OnSend(data)
}

// SendTo sends a unitcast/multicast message to all sockets in a room
func (sc *SocketClient) SendTo(roomName string, data map[string]interface{}) {
	data["Gm_Ws_Type"] = "gm_ws_multicast"
	data["Gm_Ws_Room"] = roomName

	err := websocket.JSON.Send(sc.Conn, data)
	if err != nil {
		sc.events.OnSendError(err)
		return
	}
	sc.events.OnSend(data)
}







