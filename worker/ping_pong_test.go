package worker

import (
	"fmt"
	"testing"
	"time"
)

var m = NewManager()

type ping struct{}
type pong struct{}
type pinger struct{}
type ponger struct{}

func (p *pinger) Receive(msg any) {
	switch msg.(type) {
	case pong:
		fmt.Println("ping!")
		time.Sleep(time.Millisecond * 200)
		m.Send(ping{})
	}
}

func (p *ponger) Receive(msg any) {
	switch msg.(type) {
	case ping:
		fmt.Println("pong!")
		time.Sleep(time.Millisecond * 200)
		m.Send(pong{})
	}
}

func TestPingPong(t *testing.T) {
	m.Add(&pinger{}, Metadata{})
	m.Add(&ponger{}, Metadata{})
	m.Start()
	m.Send(ping{})
	defer m.Stop()
	time.Sleep(time.Second * 15)
}
