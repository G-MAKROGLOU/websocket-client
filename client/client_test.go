package client

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/G-MAKROGLOU/websocket-server/server"
	"github.com/stretchr/testify/assert"
)


func TestConnect(t *testing.T){
	var waitG sync.WaitGroup
	var s *server.SocketServer

	waitG.Add(1)


	go func() {
		s = server.New(server.NOOPSocketServerEvents{})
		err := s.Start();
		assert.Equal(t, errors.New("http: Server closed"), err)			
	}()

	go func() {
		time.Sleep(15 * time.Second)
		s.Stop()
		waitG.Done()
	}()


	time.Sleep(5 * time.Second)

	c := New("http://localhost:3000", "ws://localhost:3000/ws", NOOPSocketClientEvents{})

	err := c.Connect()

	assert.Nil(t, err)

	assert.NotNil(t, c.Conn)

	time.Sleep(5 * time.Second)

	if c.Conn != nil {
		err := c.Disconnect()
		assert.Nil(t, nil, err)
	}

	waitG.Wait()
}
