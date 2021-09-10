package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var profAttr = [10]string{"professor_name", "quality_score", "difficulty_score", "sentiment_score_discrete", "sentiment_score_continuous", "keywords", "would_take_again", "pid", "school", "department"}
var profPIDAttr = [4]string{"professor_name", "department", "school", "pid"}

// Check and Try to retrieve query entry from Redis
func RedisCheckResultCache(redisClient *redis.Client, input string, ctx context.Context) map[string]string {
	// check Redis (as a fast RAM-based cache) if the input query exists
	res, err := redisClient.HGetAll(ctx, input).Result()
	if err != nil {
		log.Fatal("CheckCache err:", err)
	}

	return res
}

// Update Redis entry -> {PID: results} in form of HashMap
func RedisUpdateResultCache(redisClient *redis.Client, input string, values []string, ctx context.Context) int64 {
	m := make(map[string]interface{})
	// populate the map m with the new values
	for ind, val := range profAttr {
		m[val] = values[ind]
	}
	// store the query result into Redis
	ok, err := redisClient.HSet(ctx, input, m).Result()
	if err != nil {
		log.Fatal("UpdateCache err:", err)
	}

	// set the expiration time of cache --> random in range of 1 to 6 hours
	rand.Seed(time.Now().Unix())
	_, expErr := redisClient.Expire(ctx, input, time.Duration((rand.Intn(6)+1))*60*60*time.Second).Result()
	if expErr != nil {
		log.Fatal("SetExpire err:", expErr)
	}

	return ok
}

// Check and Try to retrieve query entry from Redis
func RedisCheckNameCache(redisClient *redis.Client, input string, ctx context.Context) map[string]string {
	// check Redis (as a fast RAM-based cache) if the input query exists
	res, err := redisClient.HGetAll(ctx, input).Result()
	if err != nil {
		log.Fatal("CheckCache err:", err)
	}

	return res
}

// Update Redis entry -> {Prof name: results} in form of Key/Value pair
func RedisUpdateNameCache(redisClient *redis.Client, input string, values []string, ctx context.Context) int64 {
	m := make(map[string]interface{})
	// populate the map m with the new values
	for ind, val := range values {
		m[strconv.Itoa(ind)] = val
	}

	// store the query result into Redis
	ok, err := redisClient.HSet(ctx, input, m).Result()
	if err != nil {
		log.Fatal("UpdateCache err:", err)
	}

	// set the expiration time of cache --> random in range of 1 to 6 hours
	rand.Seed(time.Now().Unix())
	_, expErr := redisClient.Expire(ctx, input, time.Duration((rand.Intn(6)+1))*60*60*time.Second).Result()
	if expErr != nil {
		log.Fatal("SetExpire err:", expErr)
	}

	return ok
}

// Check and Try to retrieve query entry from Redis
func RedisCheckDepartmentList(redisClient *redis.Client, input string, ctx context.Context) map[string]string {
	// check Redis (as a fast RAM-based cache) if the input query exists
	res, err := redisClient.HGetAll(ctx, input).Result()
	if err != nil {
		log.Fatal("CheckCache err:", err)
	}

	return res
}

// Update Redis entry -> {School name: departments list} in form of Key/Value pair
func RedisUpdateDepartmentList(redisClient *redis.Client, input string, values []string, ctx context.Context) int64 {
	m := make(map[string]interface{})
	// populate the map m with the new values
	for ind, val := range values {
		m[strconv.Itoa(ind)] = val
	}

	// store the query result into Redis
	ok, err := redisClient.HSet(ctx, input, m).Result()
	if err != nil {
		log.Fatal("UpdateCache err:", err)
	}

	// no expiration time because there is no way to retrive RMP's department data JIT...
	/*
		// set the expiration time of cache --> random in range of 1 to 6 hours
		rand.Seed(time.Now().Unix())
		_, expErr := redisClient.Expire(ctx, input, time.Duration((rand.Intn(6) + 1)) * 60 * 60 * time.Second).Result()
		if expErr != nil {
			log.Fatal("SetExpire err:", expErr)
		}
	*/

	return ok
}
