package test

import (
	"net"
	"own-redis/protocol"
	"reflect"
	"testing"
)

// 发送Redis命令
func sendCommand(conn net.Conn, cmd [][]byte) {
	var request []protocol.Reply
	for _, arg := range cmd {
		request = append(request, &protocol.BulkStringReply{Str: arg})
	}
	conn.Write((&protocol.ArrayReply{Val: request}).Encode())
}

// 读取服务器回复
func readReply(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func TestRPush(t *testing.T) {
	// 测试RPUSH命令
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// 发送RPUSH命令
	cmd := [][]byte{[]byte("RPUSH"), []byte("mykey"), []byte("value1"), []byte("value2")}
	sendCommand(conn, cmd)

	// 检查返回值
	reply, err := readReply(conn)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}
	expected := &protocol.IntegerReply{Val: 2}
	if !reflect.DeepEqual(reply, expected.Encode()) {
		t.Errorf("RPUSH command failed. Expected %v, got %v", expected, reply)
	}

	nextcmd := [][]byte{[]byte("RPUSH"), []byte("mykey"), []byte("value3")}
	sendCommand(conn, nextcmd)

	// 检查返回值
	reply, err = readReply(conn)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}
	expected = &protocol.IntegerReply{Val: 3}
	if !reflect.DeepEqual(reply, expected.Encode()) {
		t.Errorf("RPUSH command failed. Expected %v, got %v", expected, reply)
	}
}

func TestLPush(t *testing.T) {
	// 测试LPUSH命令
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// 发送LPUSH命令
	cmd := [][]byte{[]byte("LPUSH"), []byte("mykey"), []byte("value1"), []byte("value2")}
	sendCommand(conn, cmd)

	// 检查返回值
	reply, err := readReply(conn)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}
	expected := &protocol.IntegerReply{Val: 2}
	if !reflect.DeepEqual(reply, expected.Encode()) {
		t.Errorf("LPUSH command failed. Expected %v, got %v", expected, reply)
	}

	nextcmd := [][]byte{[]byte("LPUSH"), []byte("mykey"), []byte("value3")}
	sendCommand(conn, nextcmd)

	// 检查返回值
	reply, err = readReply(conn)
	if err != nil {
		t.Fatalf("Error reading reply: %v", err)
	}
	expected = &protocol.IntegerReply{Val: 3}
	if !reflect.DeepEqual(reply, expected.Encode()) {
		t.Errorf("LPUSH command failed. Expected %v, got %v", expected, reply)
	}

}
