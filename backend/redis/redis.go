package main

import (
	"errors"
	"github.com/go-redis/redis"
)

type RedisMetadata struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Redis struct {
	client *redis.Client
}

func NewRedisClient(url string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	r := &Redis{
		client: client,
	}
	return r, nil
}
func (r *Redis) Set(key string, value []byte) error {
	err := r.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	result, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	}
	if err != nil {
		return nil, err
	}
	return []byte(result), nil

}
