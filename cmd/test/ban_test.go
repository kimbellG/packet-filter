//go:build integration
// +build integration

package test

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestBanPackage(t *testing.T) {
	resp, err := http.Get("http://192.0.2.1:8080/filter/TCP")
	if err != nil {
		t.Fatalf("Failed to http-server connect: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code of http connect is %v", resp.StatusCode)
	}

	if serverHears(t) {
		t.Errorf("data came to server")
	}
}

func serverHears(t *testing.T) bool {
	valueString := "test string"
	var err error

	_, err = cl.Write([]byte(valueString))
	if err != nil {
		t.Fatalf("Failed sending bytes to client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ret := make(chan struct{})

	go func() {
		buffer := make([]byte, len(valueString)+1)

		_, err := srv.Read(buffer)
		if err != nil {
			t.Fatalf("Failed read from server connection: %v", err)
		}

		ret <- struct{}{}

	}()

	select {
	case <-ctx.Done():
		return false
	case <-ret:
		return true
	}
}

func TestUnbanPackage(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://192.0.2.1:8080/filter/TCP", nil)
	if err != nil {
		t.Fatalf("Failed to create unblock request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to connect to the server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code of http connect is %v", resp.StatusCode)
	}

	if !serverHears(t) {
		t.Errorf("data didn't come to server")
	}

}
