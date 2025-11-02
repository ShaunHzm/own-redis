package main

import (
	"fmt"
	"net"

	"own-redis/command"
)

func main() {

	fmt.Println("Logs for own-redis")

	l, err := net.Listen("tcp", "localhost:6379") // 监听 localhost:6379
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	defer l.Close() // 关闭监听

	// 不断接受连接
	for {
		conn, err := l.Accept() // 接受连接
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}
		// 使用command包中的HandleConnection函数处理连接
		// 注意：这里不需要defer conn.Close()，因为HandleConnection内部已经处理了
		go command.HandleConnection(conn)
	}

}
