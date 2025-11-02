package data

import (
	"sync"
	"time"
)

// Data 结构体表示键值对
type Data struct {
	Key   []byte
	Value []byte
}

type ExpireData struct {
	Value  string
	Expire int64
}

// KVSTORE 是一个全局的键值存储，使用 sync.Map 来确保并发安全
var (
	KVSTORE sync.Map
)

func getExpireData(key string) (ExpireData, bool) {
	raw, ok := KVSTORE.Load(key)
	if !ok {
		return ExpireData{}, false
	}
	return raw.(ExpireData), true
}

func GetValue(key string) (string, bool) {
	expireData, ok := getExpireData(key)
	if expireData.Expire > 0 && expireData.Expire < time.Now().UnixMilli() {
		return "", false // 数据过期
	}
	return expireData.Value, ok
}
