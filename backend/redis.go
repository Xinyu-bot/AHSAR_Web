package main

import (
	"log"
	"fmt"
	"github.com/go-redis/redis"
)

// Check and Try to retrieve query entry from Redis
func RedisCheckCache(redisClient *redis.Client, input string) map[string]string {
	// Redis Command
	// check Redis (as a fast RAM-based cache) if the input query exists
	res, err := redisClient.HGetAll(input).Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("res in redis.go/RedisCheckCache:", res)

	return res
}

// Update Redis entry
func RedisUpdateCache(redisClient *redis.Client, input string, values []string) string {
	m := make(map[string]interface{})
	/*
	populate the map m with the new values
	profAttr defined in app.go: [6]string{...}
	*/
	for ind, val := range profAttr {
		m[val] = values[ind]
	}
	// Redis Command
	// store the query result into Redis
	ok, err := redisClient.HMSet(input, m).Result()
	if err != nil {
		log.Println(err)
	}
	
	return ok
}