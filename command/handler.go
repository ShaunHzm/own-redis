package command

import (
	"fmt"
	"io"
	"net"
	"own-redis/protocol"
	"strings"
)

// HandleConnection 处理客户端连接
// 这个函数接收一个TCP连接，并负责处理客户端发送的Redis命令
func HandleConnection(conn net.Conn) {
	// 打印连接信息
	fmt.Println("New connection established from:", conn.RemoteAddr())

	// 确保连接最终被关闭
	defer conn.Close()

	// 持续处理客户端请求
	for {
		// 获取客户端请求
		req, err := protocol.Decode(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected:", conn.RemoteAddr())
				return
			}
			fmt.Println("Error decoding request:", err)
			conn.Write((&protocol.ErrorReply{Str: []byte("ERR unknown command")}).Encode())
			continue // 修改为continue而不是return
		}

		switch req := req.(type) {
		// ping 命令通过简单字符串 "PING" 来触发
		case protocol.SimpleStringRequest:
			if strings.ToUpper(string(req.Str)) == "PING" {
				fn, ok := commandMap["PING"]
				if !ok {
					fmt.Println("Error: command not found")
					conn.Write((&protocol.ErrorReply{Str: []byte("ERR unknown command")}).Encode())
					continue
				}
				reply := fn([]protocol.Request{})
				conn.Write(reply.Encode())
			} else {
				conn.Write((&protocol.ErrorReply{Str: []byte("ERR unknown command")}).Encode())
			}
		// set, get, del 命令通过数组请求触发
		case protocol.ArrayRequest:
			if len(req.Val) == 0 {
				conn.Write((&protocol.ErrorReply{Str: []byte("ERR empty command")}).Encode())
				continue
			}
			// 类型断言确保第一个元素是BulkStringRequest
			firstArg, ok := req.Val[0].(protocol.BulkStringRequest)
			if !ok {
				conn.Write((&protocol.ErrorReply{Str: []byte("ERR invalid command format")}).Encode())
				continue
			}
			cmd := strings.ToUpper(string(firstArg.Str))
			fn, ok := commandMap[cmd]
			if !ok {
				fmt.Println("Error: command not found")
				conn.Write((&protocol.ErrorReply{Str: []byte("ERR unknown command")}).Encode())
				continue
			}
			reply := fn(req.Val[1:]) // 执行对应的命令回调
			if reply == nil {
				conn.Write((&protocol.ErrorReply{Str: []byte("ERR command not implemented")}).Encode())
				continue
			}
			conn.Write(reply.Encode())
		// 其他未知请求类型
		default:
			conn.Write((&protocol.ErrorReply{Str: []byte("ERR unknown command")}).Encode())
		}
	}
}
