package command

import (
	"fmt"
	"own-redis/data"
	"own-redis/protocol"
	"strconv"
	"strings"
	"time"
)

type CommandFunc func(args []protocol.Request) protocol.Reply

var (
	commandMap = make(map[string]CommandFunc)
)

func init() {
	commandMap["PING"] = Ping
	commandMap["ECHO"] = Echo
	commandMap["SET"] = Set
	commandMap["GET"] = Get
	commandMap["DEL"] = Del
}

func Ping(args []protocol.Request) protocol.Reply {
	if len(args) > 0 {
		return &protocol.ErrorReply{Str: []byte("ERR wrong number of arguments for 'ping' command")}
	}
	return &protocol.SimpleStringReply{Str: []byte("PONG")}
}

func Echo(args []protocol.Request) protocol.Reply {
	// echo 命令需要一个参数
	if len(args) != 1 {
		return &protocol.ErrorReply{Str: []byte("ERR wrong number of arguments for 'echo' command")}
	}
	return &protocol.BulkStringReply{Str: args[0].(protocol.BulkStringRequest).Str}
}

func Set(args []protocol.Request) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ErrorReply{Str: []byte("ERR 'set' command requires at least 2 arguments")}
	}

	key := string(args[0].(protocol.BulkStringRequest).Str)
	value := string(args[1].(protocol.BulkStringRequest).Str)

	var expire int64 = 0

	// 检查是否有过期时间参数
	for i := 2; i < len(args); i++ {
		req := args[i]
		switch strings.ToUpper(string(req.(protocol.BulkStringRequest).Str)) {
		case "EX":
			if i+1 >= len(args) {
				return &protocol.ErrorReply{Str: []byte("ERR 'set' command requires an expiration time after 'ex'")}
			}
			sec, _ := strconv.ParseInt(string(args[i+1].(protocol.BulkStringRequest).Str), 10, 64)
			if expire != 0 {
				return &protocol.ErrorReply{Str: []byte("ERR 'set' command requires only one expiration time parameter")}
			}
			expire = time.Now().UnixMilli() + sec*1000
			i++ // 跳过过期时间值
		case "PX":
			if i+1 >= len(args) {
				return &protocol.ErrorReply{Str: []byte("ERR 'set' command requires an expiration time after 'px'")}
			}
			msec, _ := strconv.ParseInt(string(args[i+1].(protocol.BulkStringRequest).Str), 10, 64)
			if expire != 0 {
				return &protocol.ErrorReply{Str: []byte("ERR 'set' command requires only one expiration time parameter")}
			}
			expire = time.Now().UnixMilli() + msec
			i++ // 跳过过期时间值
		default:
			return &protocol.ErrorReply{Str: []byte("syntax error")}
		}
	}

	fmt.Printf("key: %s, value: %s, expire: %d\n", key, value, expire)
	// 存储键值对(key->ExpireData)
	data.KVSTORE.Store(string(key), data.ExpireData{Value: value, Expire: expire})

	return &protocol.SimpleStringReply{Str: []byte("OK")}
}

func Get(args []protocol.Request) protocol.Reply {
	if len(args) != 1 {
		return &protocol.ErrorReply{Str: []byte("ERR wrong number of arguments for 'get' command")}
	}
	key := string(args[0].(protocol.BulkStringRequest).Str)
	// 从KVSTORE中获取键对应的值
	value, ok := data.GetValue(key)
	if !ok {
		return &protocol.ErrorReply{Str: []byte("ERR key not found")}
	}
	return &protocol.BulkStringReply{Str: []byte(value)}
}

func Del(args []protocol.Request) protocol.Reply {
	return nil
}
