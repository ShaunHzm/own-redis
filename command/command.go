package command

import (
	"own-redis/protocol"
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

	return nil
}

func Get(args []protocol.Request) protocol.Reply {
	return nil
}

func Del(args []protocol.Request) protocol.Reply {
	return nil
}
