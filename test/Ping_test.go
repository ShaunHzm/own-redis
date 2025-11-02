package test

import (
	"fmt"
	"net"
	"testing"
)

var PINGCMD = []byte("*1\r\n$4\r\nPING\r\n")

func TestPing_own_redis(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	conn.Write(PINGCMD)

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}

	expectedReply := []byte("+PONG\r\n")
	fmt.Printf("Expected bytes: %q\n", expectedReply)
	fmt.Printf("Actual bytes: %q\n", reply[:n])
	if string(reply[:n]) != string(expectedReply) {
		t.Errorf("TestPing_own_redis: Expected reply %s, got %s", expectedReply, reply[:n])
	}
}
