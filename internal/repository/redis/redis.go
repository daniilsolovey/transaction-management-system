package redis

import (
	"context"
)

const (
	ClicksDefaultName = "clicks"
)

func (r *Repository) InsertClickKey(ctx context.Context, key string) error {
	err := r.client.Incr(key).Err()
	if err != nil {
		r.log.Error("Redis INCR failed", "err", err, "key", key)
	}
	return err
}

func (r *Repository) GetKeys(ctx context.Context) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)

	pattern := ClicksDefaultName + ":*"

	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = r.client.Scan(cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (r *Repository) GetValueByKey(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(key).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *Repository) Del(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}
