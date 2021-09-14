//go:build integration
// +build integration

package test

import (
	"context"
	_ "embed"
	"log"
	"net"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/kimbellG/packet-filter/service"
	"github.com/vishvananda/netns"
)

//go:embed xdp.source
var source string

var (
	cl  net.Conn
	srv net.Conn
)

func TestMain(m *testing.M) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		service.StartService(ctx, source)
	}()

	time.Sleep(time.Second * 3)

	listener, err := net.Listen("tcp", "192.0.2.1:8000")
	if err != nil {
		log.Fatalf("Failed to initialize socket: %v\n", err)
	}
	defer listener.Close()

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		_, err := netns.Get()
		if err != nil {
			log.Fatalf("Failed to get current netns: %v", err)
		}

		xdpRemote, err := netns.GetFromName(os.Getenv("TESTNETNS"))
		if err != nil {
			log.Fatalf("Failed to get netns: %v", err)
		}

		if err := netns.Set(xdpRemote); err != nil {
			log.Fatalf("Failed to set test netns: %v", err)
		}

		cl, err = net.Dial("tcp", "192.0.2.1:8000")
		if err != nil {
			log.Fatalf("Failed to connecting server: %v\n", err)
		}
	}()

	srv, err = listener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection: %v\n", err)
	}
	defer srv.Close()

	m.Run()
	cl.Close()
}
