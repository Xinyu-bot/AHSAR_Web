package main

import (
	"log"
	"fmt"
	"github.com/go-redis/redis"
)

func RedisCheckCache(redisClient *redis.Client, input string) map[string]string {
	// check redis (as a fast RAM-based cache) if the input query exists
	res, err := redisClient.HGetAll(input).Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("res in redis.go/RedisCheckCache:", res)

	return res
}

func RedisUpdateCache(redisClient *redis.Client, input string, values []string) string {
	m := make(map[string]interface{})
	m["first_name"] = values[0]
	m["last_name"] = values[1]
	m["quality_score"] = values[2]
	m["difficulty_score"] = values[3]
	m["sentiment_score_discrete"] = values[4]
	m["sentiment_score_continuous"] = values[5]
	ok, err := redisClient.HMSet(input, m).Result()
	if err != nil {
		log.Println(err)
	}
	
	return ok
}