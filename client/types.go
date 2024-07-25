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
	onConnect(ws *websocket.Conn, sessID string)
	onReceive(data map[string]interface{})
	onReceiveError(err error)
	onJoin(roomName string)
	onJoinError(roomName string, err error)
	onLeave(roomName string)
	onLeaveError(roomName string, err error)
	onSend(data map[string]interface{})
	onSendError(err error)
}

// NOOPSocketClientEvents is a default struct that has no implementation for the Server events
type NOOPSocketClientEvents struct {}


func (N NOOPSocketClientEvents) onConnect(ws *websocket.Conn, sessID string){}
func (N NOOPSocketClientEvents) onDisconnect(){}
func (N NOOPSocketClientEvents) onDisconnectError(err error){}
func (N NOOPSocketClientEvents) onReceive(data map[string]interface{}){}
func (N NOOPSocketClientEvents) onReceiveError(err error){}
func (N NOOPSocketClientEvents) onJoin(roomName string){}
func (N NOOPSocketClientEvents) onJoinError(roomName string, err error){}
func (N NOOPSocketClientEvents) onLeave(roomName string){}
func (N NOOPSocketClientEvents) onLeaveError(roomName string, err error){}
func (N NOOPSocketClientEvents) onSend(data map[string]interface{}){}
func (N NOOPSocketClientEvents) onSendError(err error){}
