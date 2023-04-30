package main

import (
	"fmt"
	"time"
)

func main() {
	redisAddrs := []string{"localhost:6379", "localhost:6380", "localhost:6381"}

	quorum := len(redisAddrs)/2 + 1

	redlock := NewRedlock(redisAddrs, quorum)

	lockKey := "myLock"

	ttl := time.Second * 10

	value, err := redlock.Acquire(lockKey, ttl)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Lock acquired with value:", value)

	err = redlock.Release(lockKey, value)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Lock Released!")

}
