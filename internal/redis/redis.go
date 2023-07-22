package redis

import (
	"sync"
	"time"
)

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
