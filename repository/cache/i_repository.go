package cache

import (
	"time"
)

type IRepository interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Expire(key string, duration time.Duration) error
}
