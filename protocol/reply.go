package protocol

import "fmt"

// Reply 表示 Redis 回复
// 所有 Redis 回复都必须实现 Encode 方法
// Encode 方法返回 Redis 回复的字符串表示
type Reply interface {
	Encode() []byte
}

type SimpleStringReply struct {
	Str []byte // 简单字符串回复
}

type ErrorReply struct {
	Str []byte // 错误回复
}

type IntegerReply struct {
	Val int64 // 整数回复
}

type BulkStringReply struct {
	Str []byte // 批量字符串回复
}

type ArrayReply struct {
	Val []Reply // 数组回复
}

// Encode 方法返回 Redis 回复的字符串表示
func (r *SimpleStringReply) Encode() []byte {
	return []byte(fmt.Sprintf("+%s\r\n", string(r.Str)))
}

func (r *ErrorReply) Encode() []byte {
	return []byte(fmt.Sprintf("-%s\r\n", string(r.Str)))
}

func (r *IntegerReply) Encode() []byte {
	return []byte(fmt.Sprintf(":%d\r\n", r.Val))
}

func (r *BulkStringReply) Encode() []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(r.Str), r.Str))
}

func (r *ArrayReply) Encode() []byte {
	str := []byte(fmt.Sprintf("*%d\r\n", len(r.Val)))

	for _, elem := range r.Val {
		str = append(str, elem.Encode()...)
	}

	return str
}
