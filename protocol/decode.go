package protocol

import (
	"bufio"
	"fmt"
	"io"
)

type Request interface {
}

type SimpleStringRequest struct {
	Str []byte // 简单字符串请求
}

type ErrorStringRequest struct {
	Str []byte // 错误字符串请求
}

type IntegerRequest struct {
	val int64 // 整数请求
}

type BulkStringRequest struct {
	Str []byte // 批量字符串请求
}

// ArrayRequest 数组请求
type ArrayRequest struct {
	Val []Request // 数组请求
}

// Decode 解码 RESP 协议
func Decode(conn io.Reader) (Request, error) {
	// 读取第一个字节，判断数据类型
	br := bufio.NewReader(conn)
	buf := make([]byte, 1)
	if _, err := br.Read(buf); err != nil {
		return nil, err
	}
	switch buf[0] {
	case '+':
		return decodeSimpleString(br)
	case '-':
		return decodeError(br)
	case ':':
		return decodeInteger(br)
	case '$':
		return decodeBulkString(br)
	case '*':
		return decodeArray(br)
	default:
		return nil, fmt.Errorf("unknown data type: %v", buf[0])
	}
}

// readline 读取一行数据, 检查是否以 \r\n 结尾，移除换行符，返回字符串
func readline(br *bufio.Reader) ([]byte, error) {
	buf, err := br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(buf) < 2 || buf[len(buf)-2] != '\r' {
		return nil, fmt.Errorf("line must end with \\r\\n: %v", buf)
	}
	// 移除换行符
	buf = buf[:len(buf)-2]
	if len(buf) == 0 {
		return nil, fmt.Errorf("line must not be empty")
	}
	return buf, nil
}

// decodeSimpleString 解码简单字符串
func decodeSimpleString(br *bufio.Reader) (SimpleStringRequest, error) {
	buf, err := readline(br)
	if err != nil {
		return SimpleStringRequest{}, err
	}

	return SimpleStringRequest{Str: buf}, nil
}

// decodeError 解码错误
func decodeError(br *bufio.Reader) (ErrorStringRequest, error) {
	buf, err := readline(br)
	if err != nil {
		return ErrorStringRequest{}, err
	}

	return ErrorStringRequest{Str: buf}, nil
}

// decodeInteger 解码整数
func decodeInteger(br *bufio.Reader) (IntegerRequest, error) {
	buf, err := readline(br)
	if err != nil {
		return IntegerRequest{}, err
	}

	var i int64
	_, err = fmt.Sscanf(string(buf), "%d", &i)
	if err != nil {
		return IntegerRequest{}, fmt.Errorf("invalid integer format: %v", err)
	}

	return IntegerRequest{val: i}, nil
}

// decodeBulkString 解码批量字符串
// 先解析字符串长度
// 再解析字符串
func decodeBulkString(br *bufio.Reader) (BulkStringRequest, error) {
	// 读取字符串长度
	lenStr, err := readline(br)
	if err != nil {
		return BulkStringRequest{}, err
	}
	var l int64 // l 字符串长度
	_, err = fmt.Sscanf(string(lenStr), "%d", &l)
	if err != nil {
		return BulkStringRequest{}, fmt.Errorf("invalid bulk string length format: %v", err)
	}

	// 读取字符串数据
	buf, err := readline(br)
	if err != nil {
		return BulkStringRequest{}, err
	}
	if len(buf) != int(l) {
		return BulkStringRequest{}, fmt.Errorf("bulk string length mismatch: expected %d, got %d", l, len(buf))
	}

	return BulkStringRequest{Str: buf}, nil
}

func decodeArray(br *bufio.Reader) (ArrayRequest, error) {
	lenArr, err := readline(br)
	if err != nil {
		return ArrayRequest{}, err
	}
	var l int64 // l 数组长度
	_, err = fmt.Sscanf(string(lenArr), "%d", &l)
	if err != nil {
		return ArrayRequest{}, fmt.Errorf("invalid array length format: %v", err)
	}

	// 读取数组元素
	var arr []Request
	for i := 0; i < int(l); i++ {
		elem, err := Decode(br) // 递归解码数组元素
		if err != nil {
			return ArrayRequest{}, err
		}
		arr = append(arr, elem)
	}
	return ArrayRequest{Val: arr}, nil
}
