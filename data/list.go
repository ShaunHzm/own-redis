package data

import (
	"own-redis/protocol"
	"sync"
)

type ListEntry struct {
	mu    sync.RWMutex
	value []string
}

func (l *ListEntry) RPush(value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = append(l.value, value)
}

func (l *ListEntry) LPush(value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = append([]string{value}, l.value...)
}

func (l *ListEntry) Range(start, end int64) []protocol.Reply {
	l.mu.RLock()
	defer l.mu.RUnlock()
	len := l.Length()
	if start < 0 {
		start = len + start
	}
	if end < 0 {
		end = len + end
	}
	if start < 0 {
		start = 0
	}
	if end >= len {
		end = len - 1
	}

	replies := make([]protocol.Reply, end-start+1)
	for i := start; i <= end; i++ {
		replies[i-start] = &protocol.BulkStringReply{Str: []byte(l.value[i])}
	}
	return replies

}

func (l *ListEntry) Length() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return int64(len(l.value))
}
