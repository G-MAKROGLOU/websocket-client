# INSTALLATION 

```go
go get github.com/G-MAKROGLOU/websocket-client@latest

```


# JSON client for connection with a server created with:

```github.com/G-MAKROGLOU/websocket-server```

## CLIENT (clientevents.go)

<p>
    Implement the events interface with the functionality you want to be implemented per event. 
    Leave empty for no actions on a specific event.
</p>


```go
package main

import (
    "encoding/json"
    "fmt"
)

// Events implements SocketClientEvents interface
type Events struct {}

// OnDisconnectError event interface implementation
func (e Events) OnDisconnectError(err error){
    fmt.Println("OnDisconnectError: ", err)
}

// OnReceive event interface implementation
func (e Events) OnReceive(data map[string]interface{}){
    fmt.Println("OnReceive: ", data)
}

// OnReceiveError event interface implementation
func (e Events) OnReceiveError(err error){
    if err == io.EOF {
        fmt.Println("connection closed ")
    } else {
        fmt.Println("some error during receive")
    }
}

// OnJoinError event interface implementation
func (e Events) OnJoinError(roomName string, err error){
    fmt.Println("OnJoinError: ", roomName, err)
}

// OnLeaveError event interface implementation
func (e Events) OnLeaveError(roomName string, err error){
    fmt.Println("OnLeaveError: ", roomName, err)
}

// OnSendError event interface implementation
func (e Events) OnSendError(err error){
    if err == err.(*net.OpError) {
	slog.Error("server was closed")
    } else {
	slog.Error("unexpected send error: ", "reason", err)
    }
}

```

## CLIENT (client.go)

<p>
    Connect and interact with a socket server made with github.com/G-MAKROGLOU/websocket-server
</p>


```go
package main

import (
    "log/slog"
    client "github.com/G-MAKROGLOU/websocket-client"
)

func main() {
    
    // setup the client 
    origin := "http://localhost:3000"
    server := "ws://localhost:3000/ws"

    //connect
    c1 := client.New(origin, server, Events{})
    if err := c1.Connect(); err != nil {
        slog.Error("failed to connect to server")
        os.Exit(1)
    }

    // start listening for incoming messages in a goroutine
    go c1.ReceiveJSON()

    // broadcast a message to all connected clients
    data = map[string]interface{}{
	"msg": "pong",
    }
    c1.SendJSON(data)

    // join a room
    c1.Join("test")

    // mutlicast a message to a room
    c1.SendJSONTo("test", data)

    // leave a room
    c1.Leave("test")

    // disconnect
    c1.Disconnect()
}

```


# Plain Text client for connection with any server that sends text (e.g aisstream)

## CLIENT (clientevents.go)

<p>
    Implement the events interface with the functionality you want to be implemented per event. 
    Leave empty for no actions on a specific event.
</p>


```go
package main

import (
    "encoding/json"
    "fmt"
    aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
)

// Events implements SocketClientEvents interface
type Events struct {}

// OnDisconnectError event interface implementation
func (e Events) OnDisconnectError(err error){
    fmt.Println("OnDisconnectError: ", err)
}

// OnReceive event interface implementation
func (e Events) OnReceive(data map[string]interface{}){
    var res aisstream.AisStreamMessage

    b, _ := json.Marshal(data)

    json.Unmarshal(b, &res)

    switch(res.MessageType) {
    case aisstream.POSITION_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.UNKNOWN_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.ADDRESSED_SAFETY_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.ADDRESSED_BINARY_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.AIDS_TO_NAVIGATION_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.ASSIGNED_MODE_COMMAND:
	fmt.Println(res.MessageType)
	break
    case aisstream.BASE_STATION_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.BINARY_ACKNOWLEDGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.BINARY_BROADCAST_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.CHANNEL_MANAGEMENT:
	fmt.Println(res.MessageType)
	break
    case aisstream.COORDINATED_UTC_INQUIRY:
	fmt.Println(res.MessageType)
	break
    case aisstream.DATA_LINK_MANAGEMENT_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.DATA_LINK_MANAGEMENT_MESSAGE_DATA:
	fmt.Println(res.MessageType)
	break
    case aisstream.EXTENDED_CLASS_B_POSITION_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.GROUP_ASSIGNMENT_COMMAND:
	fmt.Println(res.MessageType)
	break
    case aisstream.GNSS_BROADCAST_BINARY_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.INTERROGATION:
	fmt.Println(res.MessageType)
	break
    case aisstream.LONG_RANGE_AIS_BROADCAST_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.MULTI_SLOT_BINARY_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.SAFETY_BROADCAST_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.SHIP_STATIC_DATA:
	fmt.Println(res.MessageType)
	break
    case aisstream.SINGLE_SLOT_BINARY_MESSAGE:
	fmt.Println(res.MessageType)
	break
    case aisstream.STANDARD_CLASS_B_POSITION_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.STANDARD_SEARCH_AND_RESCUE_AIRCRAFT_REPORT:
	fmt.Println(res.MessageType)
	break
    case aisstream.STATIC_DATA_REPORT:
	fmt.Println(res.MessageType)
	break		
    }
}

// OnReceiveError event interface implementation
func (e Events) OnReceiveError(err error){
    if err == io.EOF {
	fmt.Println("connection closed ")
    } else {
	fmt.Println("some error during receive")
    }
}

// OnJoinError event interface implementation
func (e Events) OnJoinError(roomName string, err error){}

// OnLeaveError event interface implementation
func (e Events) OnLeaveError(roomName string, err error){}

// OnSendError event interface implementation
func (e Events) OnSendError(err error){
    if err == err.(*net.OpError) {
	slog.Error("server was closed")
    } else {
	slog.Error("unexpected send error: ", "reason", err)
    }
}

```


## CLIENT (client.go)


```go
package main

import (
    "sync"

    client "github.com/G-MAKROGLOU/websocket-client"
    aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
)

var waitG sync.WaitGroup

var apiKey = "YOUR-API-KEY"

func main() {
    // simulate blocking
    waitG.Add(1)

    // setup the client
    origin := "https://stream.aisstream.io"
    server := "wss://stream.aisstream.io/v0/stream"

    //connect
    c1 := client.New(origin, server, Events{})

    if err := c1.Connect(); err != nil {
	slog.Error("failed to connect to server")
	os.Exit(1)
    }

    // start listening for incoming messages in a goroutine. handle incoming messages
    // in OnReceive event
    go c.ReceiveText()

    // send the init message to start receiving text
    msg := aisstream.SubscriptionMessage{
	APIKey: apiKey,
	BoundingBoxes: [][][]float64{{ {-90.0, -180.0}, { 90.0, 180.0 }}},
	FiltersShipMMSI: []string{},
    }

    c.SendText(msg)

    // simulate end of operations
    go func() {
        time.Sleep(20 * time.Second)
        c.Disconnect()        
        waitG.Done()
    }()

    // simulate blocking
    waitG.Wait()
}

```


### NOTES

The events are soon going to change and be reducted to a few useful events. For example, OnError(err error) instead of multiple On*Error(err error).
