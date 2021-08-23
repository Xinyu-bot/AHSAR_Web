package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetProfByID(c *gin.Context) {
	// 	extract input from request query
	input := c.Query("input")
	noCache := c.Query("noCache")
	// 	initialize res for storing result from NLP Server
	var res []string
	fmt.Println("[LOG] query input: " + input + ", noCache: " + noCache)

	// if frontend enforces a no-cache-query, fetch latest data from RMP website
	if noCache == "true" {
		// 	obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
		res = ObtainProfessor(input)
		// 	update Redis with the new data
		if check := RedisCheckResultCache(redisClient, input); len(check) == 0 {
			RedisUpdateResultCache(redisClient, input, res)
		}
	} else { // use cache if available
		// 	check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckResultCache(redisClient, input)
		// 	data not cached in Redis
		if (len(redisRes) == 0) {
			// 	obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
			res = ObtainProfessor(input)
			// 	update Redis with the new data
			if check := RedisCheckResultCache(redisClient, input); len(check) == 0 {
				RedisUpdateResultCache(redisClient, input, res)
			}
		} else { //  retrieve cached data from Redis
			// 	populate res with the retrieved data
			res = make([]string, 8)
			for ind, key := range profAttr {
				res[ind] = redisRes[key]
			}
		}
	}

	// 	return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		profAttr[0]: res[0], profAttr[1]: res[1], profAttr[2]: res[2], 
		profAttr[3]: res[3], profAttr[4]: res[4], profAttr[5]: res[5], 
		profAttr[6]: res[6], "queryHash": hash,
	})
}