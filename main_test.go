package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	HOST = "ws://localhost"
)

type TestConfig struct {
	ClientCount int
	wg          *sync.WaitGroup
}

func DialServer(wg *sync.WaitGroup) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(fmt.Sprintf("%s%s", HOST, PORT), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		conn.Close()
		wg.Done()

	}()

	time.Sleep(2 * time.Second)
}

func TestConnection(t *testing.T) {
	go startServer()

	time.Sleep(time.Second)

	tc := TestConfig{
		ClientCount: 3,
		wg:          new(sync.WaitGroup),
	}

	tc.wg.Add(tc.ClientCount)

	for range tc.ClientCount {
		go DialServer(tc.wg)
	}

	tc.wg.Wait()

	time.Sleep(time.Second)

	log.Println("Exiting test")
}
