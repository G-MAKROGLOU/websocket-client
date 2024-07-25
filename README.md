# CLIENT (clientevents.go)

<p>
    Implement the events interface with the functionality you want to be implemented per event. 
    Leave empty for no actions on a specific event.
</p>


```go
package main

import (
    "encoding/json"
    "fmt"

    "golang.org/x/net/websocket"
)

// CustomEvents implements the SocketClientEvents interface
type CustomEvents struct {}

func (c CustomEvents) onConnect(ws *websocket.Conn, sessID string) {
    fmt.Println("[CLIENT] Connected with session ID: ", sessID)
}

func (c CustomEvents) onConnectError(err error){
    fmt.Println("[CLIENT] Failed to connect: ", err)
}

func (c CustomEvents) onDisconnect(){
    fmt.Println("[CLIENT] Disconnected" )
}

func (c CustomEvents) onDisconnectError(err error){
    fmt.Println("[CLIENT] Failed to disconnect: ", err)
}

func (c CustomEvents) onReceive(data map[string]interface{}){
    b, _ := json.MarshalIndent(data, "", " ")

    fmt.Println("[CLIENT] RECEIVED: ", string(b))
}

func (c CustomEvents) onReceiveError(err error){
    fmt.Println("[CLIENT] Failed to receive: ", err)
}

func (c CustomEvents) onJoin(roomName string){
    fmt.Println("[CLIENT] Joined room ", roomName)
}

func (c CustomEvents) onJoinError(roomName string, err error){
    fmt.Println("[CLIENT] Failed to join room: ", roomName, " ", err)
}

func (c CustomEvents) onLeave(roomName string){
    fmt.Println("[CLIENT] Left room: ", roomName)
}

func (c CustomEvents) onLeaveError(roomName string, err error){
    fmt.Println("[CLIENT] Failed to leave room: ", roomName, " ",  err)
}

func (c CustomEvents) onSend(data map[string]interface{}){
    b, _ := json.MarshalIndent(data, "", " ")

    fmt.Println("[CLIENT] SENT: ", string(b))
}

func (c CustomEvents) onSendError(err error){
    fmt.Println("[CLIENT] Failed to send: ", err)
}

```

# CLIENT (client.go)

<p>
    Start the client. Place client.Receive() in a goroutine so you have the main thread free for other operations. Free the waitGroup when you are done.
</p>


```go
package main

import (
    "fmt"
    client "github.com/G-MAKROGLOU/websocket-client"
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
