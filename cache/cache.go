package cache

import (
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found in the cache")
)

type SetParams struct {
	Key      string
	Value    string
	Duration time.Duration
}

type getter interface {
	Get(key string) (string, error)
}

type setter interface {
	Set(params *SetParams) error
}

//Cache the cache interface is responsible for
//providing basic cache operations with providers such as redis, bigcache, fastcache etc.
type Cache interface {
	getter
	setter
}
