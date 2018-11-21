package tcp

import (
	"testing"
	"fmt"
)

func TestServer_NewServer_Ok(t *testing.T) {
	testAddress := ":8080"
	testCh := make(chan []byte, 512)
	testServer := NewServer(testAddress, testCh)

	fmt.Println(testServer)
	if testServer == nil {
		t.Error("TestServer_Serve_Ok fail!")
	}
}

func TestServer_NewServer_WrongAddress(t *testing.T) {
	testAddress := " "
	testCh := make(chan []byte, 512)
	testServer := NewServer(testAddress, testCh)
	fmt.Println(testServer)

	fmt.Println(testServer)
	if testServer != nil {
		t.Error("TestServer_NewServer_WrongAddress fail!")
	}
}
