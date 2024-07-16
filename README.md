# CLIENT (clientevents.go)

<p>
    Implement the events interface with the functionality you want to be implemented per event. 
    Leave empty for no actions on a specific event.
</p>


```go
package testruns

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

// CustomEvents implements the SocketClientEvents interface
type CustomEvents struct {}

func (c CustomEvents) OnConnect(ws *websocket.Conn, sessID string) {
	fmt.Println("[CLIENT] Connected with session ID: ", sessID)
}

func (c CustomEvents) OnConnectError(err error){
	fmt.Println("[CLIENT] Failed to connect: ", err)
}

func (c CustomEvents) OnDisconnect(){
	fmt.Println("[CLIENT] Disconnected" )
}

func (c CustomEvents) OnDisconnectError(err error){
	fmt.Println("[CLIENT] Failed to disconnect: ", err)
}

func (c CustomEvents) OnReceive(data map[string]interface{}){

	b, _ := json.MarshalIndent(data, "", " ")

	fmt.Println("[CLIENT] RECEIVED: ", string(b))
}

func (c CustomEvents) OnReceiveError(err error){
	fmt.Println("[CLIENT] Failed to receive: ", err)
}

func (c CustomEvents) OnJoin(roomName string){
	fmt.Println("[CLIENT] Joined room ", roomName)
}

func (c CustomEvents) OnJoinError(roomName string, err error){
	fmt.Println("[CLIENT] Failed to join room: ", roomName, " ", err)
}

func (c CustomEvents) OnLeave(roomName string){
	fmt.Println("[CLIENT] Left room: ", roomName)
}

func (c CustomEvents) OnLeaveError(roomName string, err error){
	fmt.Println("[CLIENT] Failed to leave room: ", roomName, " ",  err)
}

func (c CustomEvents) OnSend(data map[string]interface{}){
	b, _ := json.MarshalIndent(data, "", " ")

	fmt.Println("[CLIENT] SENT: ", string(b))
}

func (c CustomEvents) OnSendError(err error){
	fmt.Println("[CLIENT] Failed to send: ", err)
}

```

# CLIENT (client.go)

<p>
    Start the client. Place client.Receive() in a goroutine so you have the main thread free for other operations.
</p>


```go
package main

import (
	"fmt"
	"github.com/G-MAKROGLOU/websocket-client"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	origin := "http://localhost"
	server := "ws://localhost:5000/ws"

	client := client.NewSocketClient(origin, server, CustomEvents{})

	client.Connect()

	client.Join("testRoom")

	go client.Receive()

	go testMulticast(client, "CLIENT1")

	wg.Wait()
}

func testMulticast(client *client.SocketClient, clientName string) {

	index := 0
	for {
		if index == 10 {
			fmt.Println("[MULTICAST] DISCONNECTING CLIENT: ", client.ID)
			client.Disconnect()
			break
		}
		time.Sleep(5 * time.Second)

		data := map[string]interface{}{
			"Message": "[FROM] [" + clientName + "] " + client.ID + " TO ROOM: testRoom",	
		}
		client.SendTo("testRoom", data)
		index++
	}
}

```
