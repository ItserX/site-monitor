package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: rdb}, nil
}

func (r *RedisStorage) AddSite(ctx context.Context, site Site) (string, error) {
	if site.ID == "" {
		site.ID = uuid.New().String()
	}

	key := fmt.Sprintf("site:%s", site.ID)
	err := r.client.HSet(ctx, key, map[string]interface{}{
		"url":    site.URL,
		"active": site.Active,
	}).Err()
	if err != nil {
		return "", err
	}

	err = r.client.SAdd(ctx, "sites", site.ID).Err()
	if err != nil {
		return "", err
	}

	return site.ID, nil
}

func (r *RedisStorage) GetSites(ctx context.Context) ([]Site, error) {
	ids, err := r.client.SMembers(ctx, "sites").Result()
	if err != nil {
		return nil, err
	}

	var sites []Site
	for _, id := range ids {
		site, err := r.GetSiteByID(ctx, id)
		if err == nil && site != nil {
			sites = append(sites, *site)
		}
	}

	return sites, nil
}

func (r *RedisStorage) GetSiteByID(ctx context.Context, id string) (*Site, error) {
	key := fmt.Sprintf("site:%s", id)
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	site := Site{
		ID:  id,
		URL: data["url"],
	}
	site.Active = data["active"] == "1" || data["active"] == "true"

	return &site, nil
}

func (r *RedisStorage) UpdateSite(ctx context.Context, site Site) error {
	key := fmt.Sprintf("site:%s", site.ID)
	return r.client.HSet(ctx, key, map[string]interface{}{
		"url":    site.URL,
		"active": site.Active,
	}).Err()
}

func (r *RedisStorage) DeleteSite(ctx context.Context, id string) error {
	key := fmt.Sprintf("site:%s", id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return r.client.SRem(ctx, "sites", id).Err()
}
