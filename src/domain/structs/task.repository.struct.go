package structs

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type TaskRepository struct {
	Client *redis.Client
	TTL    time.Duration
}
