package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Redlock struct {
	redisClients []*redis.Client
	quorum       int
}

func NewRedlock(redisAddrs []string, quorum int) *Redlock {
	clients := make([]*redis.Client, len(redisAddrs))

	for i, addr := range redisAddrs {
		clients[i] = redis.NewClient(&redis.Options{
			Addr: addr,
		})
	}

	return &Redlock{redisClients: clients, quorum: quorum}
}

func (r *Redlock) Acquire(lockKey string, ttl time.Duration) (string, error){
	var value string
	retries := 0

	for{
		value = randomString(32)
		successes := 0

		start := time.Now()

		for _, client := range r.redisClients{
			ok, err := client.SetNX(lockKey, value, ttl).Result()
			if err != nil{
				return "",err
			}
			if ok{
				successes++
			}
		}
		elapsed := time.Since(start)
		if successes >= r.quorum && elapsed<ttl{
			return value, nil
		}

		for _,client := range r.redisClients{
			if client.Get(lockKey).Val() == value{
				client.Del(lockKey)
			}
		}

		retries++
		if(retries>3){
			return "", fmt.Errorf("failed to acquire lock after %d retries", retries)
		}
		time.Sleep(randomSleep())
	}
}

func (r *Redlock) Release(lockKey string, value string) error{
	for _, client := range r.redisClients{
		if client.Get(lockKey).Val() == value{
			_, err := client.Del(lockKey).Result()
			if err!=nil{
				return err
			}
		}
	}
	return nil
}