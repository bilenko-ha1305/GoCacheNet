package commands

import (
	_ "encoding/json"
	"strings"
	"sync"
	"time"
)

type RedisCommand struct {
	Command string      `json:"command"`
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
}

func (r *Redis) HandleCommand(command string, key string, value interface{}) interface{} {
	switch strings.ToLower(command) {
	case "set":
		r.Set(key, value)
		return "OK"
	case "get":
		if val, exists := r.Get(key); exists {
			return val
		} else {
			return "Key not found"
		}
	case "del":
		r.Delete(key)
		return "OK"
	default:
		return "Unknown command"
	}
}

type Redis struct {
	data map[string]interface{}
	mu   sync.Mutex
}

func NewRedis() *Redis {
	return &Redis{
		data: make(map[string]interface{}),
	}
}

func (r *Redis) Set(key string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

func (r *Redis) Get(key string) (interface{}, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, exists := r.data[key]
	return value, exists
}

func (r *Redis) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, key)
}

func (r *Redis) Expire(key string, duration time.Duration) {
	go func() {
		<-time.After(duration)
		r.Delete(key)
	}()
}
