package command

import (
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
	commandMap["RPUSH"] = RPush
	commandMap["LPUSH"] = LPush
	commandMap["LRANGE"] = LRange
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
	if len(args) != 1 {
		return &protocol.ErrorReply{Str: []byte("ERR wrong number of arguments for 'del' command")}
	}
	key := string(args[0].(protocol.BulkStringRequest).Str)
	// 从KVSTORE中删除键对应的值
	data.KVSTORE.Delete(key)
	return &protocol.IntegerReply{Val: 1}
}

// RPush 实现 RPUSH 命令
func RPush(args []protocol.Request) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ErrorReply{Str: []byte("ERR 'rpush' command requires at least 2 arguments")}
	}
	key := string(args[0].(protocol.BulkStringRequest).Str)
	// 从KVSTORE中获取或创建列表项
	raw, _ := data.KVSTORE.LoadOrStore(key, &data.ListEntry{})
	list, ok := raw.(*data.ListEntry)
	if !ok {
		return &protocol.ErrorReply{Str: []byte("WRONGTYPE Operation against a key holding the wrong kind of value")}
	}
	for i := 1; i < len(args); i++ {
		value := string(args[i].(protocol.BulkStringRequest).Str)
		list.RPush(value)
	}
	return &protocol.IntegerReply{Val: int64(list.Length())}
}

func LPush(args []protocol.Request) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ErrorReply{Str: []byte("ERR 'lpush' command requires at least 2 arguments")}
	}
	key := string(args[0].(protocol.BulkStringRequest).Str)
	// 从KVSTORE中获取或创建列表项
	raw, _ := data.KVSTORE.LoadOrStore(key, &data.ListEntry{})
	list := raw.(*data.ListEntry)
	for i := 1; i < len(args); i++ {
		value := string(args[i].(protocol.BulkStringRequest).Str)
		list.LPush(value)
	}
	return &protocol.IntegerReply{Val: list.Length()}
}

func LRange(args []protocol.Request) protocol.Reply {
	if len(args) != 3 {
		return &protocol.ErrorReply{Str: []byte("ERR wrong number of arguments for 'lrange' command")}
	}
	key := string(args[0].(protocol.BulkStringRequest).Str)
	start, _ := strconv.ParseInt(string(args[1].(protocol.BulkStringRequest).Str), 10, 64)
	end, _ := strconv.ParseInt(string(args[2].(protocol.BulkStringRequest).Str), 10, 64)
	// 从KVSTORE中获取列表项
	raw, ok := data.KVSTORE.Load(key)
	if !ok {
		return &protocol.ArrayReply{Val: []protocol.Reply{}}
	}
	list, ok := raw.(*data.ListEntry)
	if !ok {
		return &protocol.ErrorReply{Str: []byte("WRONGTYPE Operation against a key holding the wrong kind of value")}
	}
	// 返回指定范围内的元素
	return &protocol.ArrayReply{Val: list.Range(start, end)}
}
