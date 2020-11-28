package raknet_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/sandertv/go-raknet"
)

func TestListen(t *testing.T) {
	l, err := raknet.Listen(":19132")
	if err != nil {
		panic(err)
	}
	go func() {
		_, _ = raknet.Dial("127.0.0.1:19132")
	}()
	c := make(chan error)
	go accept(l, c)

	select {
	case err := <-c:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("accepting connection took longer than 3 seconds")
	}
}

func accept(l *raknet.Listener, c chan error) {
	if _, err := l.Accept(); err != nil {
		c <- fmt.Errorf("error accepting connection: %v", err)
	}
	c <- nil
}

func TestUnconnectedPingOverride(t *testing.T) {
	const (
		overriddenPong = "overridden pong"
	)

	lc := raknet.ListenConfig{
		HandleUnconnectedPing: func(_ net.Addr) []byte {
			return []byte(overriddenPong)
		},
	}
	l, err := lc.Listen(":19132")
	if err != nil {
		panic(err)
	}
	response, err := raknet.Ping("127.0.0.1:19132")
	if err != nil {
		t.Fatalf("error connecting to server socket: %v", err)
	}
	if string(response) != overriddenPong {
		t.Fatalf("response doesn't match the overridden pong")
	}
	l.Close()
}
