package redis

import (
	"context"
)

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
