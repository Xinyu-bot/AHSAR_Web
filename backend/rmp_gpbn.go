package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetPidByName(c *gin.Context) {
	// extract input from request query
	input := c.Query("input")
	noCache := c.Query("noCache")
	// initialize res for storing result from NLP Server
	var res, ret []string
	var hasResult string
	fmt.Println("[LOG] query input: " + input + ", noCache: " + noCache)

	// if frontend enforces a no-cache-query, fetch latest data from RMP website
	if noCache == "true" {
		res = ObtainPID(input)
		if res[0] == "-1" {
			hasResult = "false"
			ret = append(ret, "DNE")
		} else {
			hasResult = "true"
			for _, profEntry := range res {
				tempEntry := strings.Split(profEntry, "&")
				ret = append(ret, fmt.Sprintf("%s %s %s %s", tempEntry[0], tempEntry[1], tempEntry[2], tempEntry[3]))
			}
		}
		// update Redis with the new data
		RedisUpdateNameCache(redisClient, input, ret, ctx)
	} else { 
		// check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckNameCache(redisClient, input, ctx)
		// data not cached in Redis
		if (len(redisRes) == 0) {
			// acquire mutex on the key in Redis
			if ok, _ := redisClient.SetNX(ctx, input + "_mutex", 1, time.Duration(5) * time.Second).Result(); ok == true {
				// obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
				res = ObtainPID(input)
				if res[0] == "-1" {
					hasResult = "false"
					ret = append(ret, "DNE")
				} else {
					hasResult = "true"
					for _, profEntry := range res {
						tempEntry := strings.Split(profEntry, "&")
						ret = append(ret, fmt.Sprintf("%s %s %s %s", tempEntry[0], tempEntry[1], tempEntry[2], tempEntry[3]))
					}
				}
				RedisUpdateNameCache(redisClient, input, ret, ctx)
				// release mutex
				redisClient.Del(ctx, input + "_mutex")
			} else {
				// fail to acquire, meaning other goroutine is holding the mutex... then sleep for 50ms
				time.Sleep(50)
				// check mutex every 50ms
				for ;; {
					checkMutex, _ := redisClient.Exists(ctx, input + "_mutex").Result()
					if checkMutex == 1 {
						time.Sleep(50)
					} else { // if mutex has been released, 
						// meaning that Redis should have cached data, then break and get the cached data
						break
					}
				}
				// obtain cached data from Redis
				redisRes = RedisCheckNameCache(redisClient, input, ctx)
				// populate res with the retrieved data
				if redisRes["0"] == "DNE" {
					hasResult = "false"
					ret = append(ret, "DNE")
				} else {
					hasResult = "true"
					for _, val := range redisRes {
						ret = append(ret, val)
					}
				}
			} 
		} else { //  retrieve cached data from Redis
			// populate res with the retrieved data
			if redisRes["0"] == "DNE" {
				hasResult = "false"
				ret = append(ret, "DNE")
			} else {
				hasResult = "true"
				fmt.Println(redisRes)
				for _, val := range redisRes {
					ret = append(ret, val)
				}
			}
		}
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		"queryHash": hash, "hasResult": hasResult, "ret": ret, 
	})
}