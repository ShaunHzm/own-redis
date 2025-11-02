package test

import (
	"fmt"
	"net"
	"testing"
)

var EXPECTEDREPLY = []byte("$10\r\n0123456789\r\n")

func TestEcho_own_redis(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	var echo_cmd = append([]byte("*2\r\n$4\r\nECHO\r\n"), EXPECTEDREPLY...)
	conn.Write(echo_cmd)

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}

	fmt.Printf("Expected bytes: %q\n", EXPECTEDREPLY)
	fmt.Printf("Actual bytes: %q\n", reply[:n])
	if string(reply[:n]) != string(EXPECTEDREPLY) {
		t.Errorf("TestEcho_own_redis: Expected reply %s, got %s", EXPECTEDREPLY, reply[:n])
	}
}
