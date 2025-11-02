package test

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func set_value(conn net.Conn, key, value []byte) {
	conn.Write([]byte("*3\r\n$3\r\nSET\r\n$" + strconv.Itoa(len(key)) + "\r\n" + string(key) + "\r\n$" + strconv.Itoa(len(value)) + "\r\n" + string(value) + "\r\n"))
}

func set_value_with_expire(conn net.Conn, key, value []byte, expire int64) {
	conn.Write([]byte("*5\r\n$3\r\nSET\r\n$" +
		strconv.Itoa(len(key)) + "\r\n" + string(key) + "\r\n$" +
		strconv.Itoa(len(value)) + "\r\n" + string(value) + "\r\n$" +
		strconv.Itoa(len("ex")) + "\r\n" + "ex" + "\r\n$" +
		strconv.Itoa(len(strconv.FormatInt(expire, 10))) + "\r\n" + strconv.FormatInt(expire, 10) + "\r\n"))
}

func get_value(conn net.Conn, key []byte) ([]byte, error) {
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$" + strconv.Itoa(len(key)) + "\r\n" + string(key) + "\r\n"))
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response[:n], nil
}

func get_value_with_expire(conn net.Conn, key []byte, expire int64) ([]byte, error) {
	time.Sleep(time.Duration(expire) * time.Second)
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$" + strconv.Itoa(len(key)) + "\r\n" + string(key) + "\r\n"))
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response[:n], nil
}

func TestSet_own_redis(t *testing.T) {
	// 测试SET命令
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// 发送SET命令
	set_value(conn, []byte("key"), []byte("value"))

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	// 检查响应是否为OK
	expected := "+OK\r\n"
	if string(reply[:n]) != expected {
		t.Errorf("Unexpected response. Expected: %s, Got: %s", expected, string(reply[:n]))
	}

	// 读取服务器响应
	response, err := get_value(conn, []byte("key"))
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	// 检查响应是否为value
	expected = "$5\r\nvalue\r\n"
	if string(response) != expected {
		t.Errorf("Unexpected response. Expected: %s, Got: %s", expected, string(response))
	}
}

func TestSetExpire_own_redis(t *testing.T) {
	// 测试SET命令
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// 发送SET命令, 设置过期时间为10秒
	set_value_with_expire(conn, []byte("key"), []byte("value"), 10)

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	// 检查响应是否为OK
	expected := "+OK\r\n"
	if string(reply[:n]) != expected {
		t.Errorf("Unexpected response. Expected: %s, Got: %s", expected, string(reply[:n]))
	}

	// 读取服务器响应, 等待过期时间
	response, err := get_value_with_expire(conn, []byte("key"), 5)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	// 检查响应是否为value
	expected = "$5\r\nvalue\r\n"
	if string(response) != expected {
		t.Errorf("Unexpected response. Expected: %s, Got: %s", expected, string(response))
	}
}
