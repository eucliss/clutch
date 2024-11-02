package receiver

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"encoding/json"

	"clutch/common"

	"github.com/gorilla/websocket"
)

type Receiver struct {
	SynthChan         *chan common.Event
	MaskedStorageChan *chan common.Event
	eventChan         *chan common.Event
	pipeline          *chan common.Event
	done              chan struct{}
	wg                sync.WaitGroup
	upgrader          websocket.Upgrader
}

func NewReceiver() *Receiver {
	return &Receiver{
		SynthChan:         &common.SynthChan,
		MaskedStorageChan: &common.MaskedStorageChan,
		eventChan:         &common.EventChan,
		pipeline:          &common.Pipeline,
		done:              make(chan struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections for this example
			},
		},
	}
}

func (r *Receiver) Receive() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case event := <-*r.eventChan:
				*r.pipeline <- event
				fmt.Printf("Received event: %+v\n", event)
			case <-r.done:
				return
			}
		}
	}()
}

func (r *Receiver) HandleWebSocket(w http.ResponseWriter, req *http.Request) {
	conn, err := r.upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Connection closed unexpectedly (HandleWebSocket): %v\n", err)
			}
			break
		}

		fmt.Println("Raw message: (HandleWebSocket)", string(message))

		var event common.Event
		decoder := json.NewDecoder(bytes.NewReader(message))
		decoder.UseNumber() // This helps preserve number precision
		err = decoder.Decode(&event)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			continue
		}

		fmt.Printf("Forwarding event to event channel (HandleWebSocket): %+v\n", event)
		*r.eventChan <- event
	}
}

func (r *Receiver) StartServer(addr string) error {
	http.HandleFunc("/ws", r.HandleWebSocket)
	return http.ListenAndServe(addr, nil)
}
