package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	userExpiration = 60 * 7
	codeExpiration = 60 * 5
)

type UserCache struct {
	rdb *redis.Client
}

func NewUserCache(rdb *redis.Client) *UserCache {
	return &UserCache{rdb: rdb}
}

func (r *UserCache) PutUser(ctx context.Context, telegram string, userJSON []byte) error {
	key := fmt.Sprintf("user:%s", telegram)
	return r.rdb.Set(ctx, key, userJSON, time.Second*userExpiration).Err()
}

func (r *UserCache) PutVerificationCode(ctx context.Context, telegram string, code int) error {
	key := fmt.Sprintf("code:%s", telegram)
	return r.rdb.Set(ctx, key, code, time.Second*codeExpiration).Err()
}

func (r *UserCache) GetUser(ctx context.Context, telegram string) ([]byte, error) {
	key := fmt.Sprintf("user:%s", telegram)
	return r.rdb.Get(ctx, key).Bytes()
}

func (r *UserCache) GetVerificationCode(ctx context.Context, telegram string) (int, error) {
	key := fmt.Sprintf("code:%s", telegram)
	return r.rdb.Get(ctx, key).Int()
}
