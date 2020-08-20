package cache

import "time"

//Storage mechanism for caching strings
type Storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
}

// goenning/go-cache
