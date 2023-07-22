package internal

import (
	_ "encoding/json"
	"strings"
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
