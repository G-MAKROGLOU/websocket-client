package client

import "golang.org/x/net/websocket"

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
	OnDisconnectError(err error)
	OnReceive(data map[string]interface{})
	OnReceiveError(err error)
	OnJoinError(roomName string, err error)
	OnLeaveError(roomName string, err error)
	OnSendError(err error)
}

// NOOPSocketClientEvents is a default struct that has no implementation for the Server events
type NOOPSocketClientEvents struct {}

// OnDisconnectError emitted when the client fails to disconnect
func (n NOOPSocketClientEvents) OnDisconnectError(err error){}

// OnReceive emitted when the client receives a message. all message handling happens here.
func (n NOOPSocketClientEvents) OnReceive(data map[string]interface{}){}

// OnReceiveError emitted when the client fails to receive a message
func (n NOOPSocketClientEvents) OnReceiveError(err error){}

// OnJoinError emitted when the client fails to join a room.
func (n NOOPSocketClientEvents) OnJoinError(roomName string, err error){}

// OnLeaveError emitted when the client fails to leave a room
func (n NOOPSocketClientEvents) OnLeaveError(roomName string, err error){}

// OnSendError emitted when the client fails to send a message
func (n NOOPSocketClientEvents) OnSendError(err error){}
