package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetPidByName(c *gin.Context) {
	// 	extract input from request query
	input := c.Query("input")
	noCache := c.Query("noCache")
	// 	initialize res for storing result from NLP Server
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
		// 	update Redis with the new data
		if check := RedisCheckNameCache(redisClient, input); len(check) == 0 {
			RedisUpdateNameCache(redisClient, input, ret)
		}
	} else { 
		// 	check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckNameCache(redisClient, input)
		// 	data not cached in Redis
		if (len(redisRes) == 0) {
			// 	obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
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
			// 	update Redis with the new data
			if check := RedisCheckNameCache(redisClient, input); len(check) == 0 {
				RedisUpdateNameCache(redisClient, input, ret)
			}
		} else { //  retrieve cached data from Redis
			// 	populate res with the retrieved data
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

	// 	return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		"queryHash": hash, "hasResult": hasResult, "ret": ret, 
	})
}