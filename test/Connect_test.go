package test

import (
	"net"
	"testing"
)

func TestConnect_own_redis(t *testing.T) {
	// 创建一个错误通道，用于从goroutine中传递错误
	errCh := make(chan error, 10)
	// 用于等待所有goroutine完成的通道
	doneCh := make(chan struct{}, 10)

	// 测试连接10次
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				doneCh <- struct{}{}
			}()

			conn, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				// 将错误发送到错误通道，而不是直接调用t.Fatalf
				errCh <- err
				return
			}
			defer conn.Close()

			// 这里可以添加更多的测试逻辑
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-doneCh
	}
	close(doneCh)

	// 检查是否有任何错误
	close(errCh)
	for err := range errCh {
		// 在主goroutine中调用t.Fatalf
		t.Fatalf("Error connecting to server: %v", err)
	}
}
