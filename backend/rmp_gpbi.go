package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
)

func GetProfByID(c *gin.Context) {
	// extract input from request query
	input := c.Query("input")
	noCache := c.Query("noCache")
	// initialize res for storing result from NLP Server
	var res []string
	fmt.Println("[LOG] query input: " + input + ", noCache: " + noCache)

	// if frontend enforces a no-cache-query, fetch latest data from RMP website
	if noCache == "true" {
		// obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
		res = ObtainProfessor(input)
		// update Redis with the new data
		RedisUpdateResultCache(redisClient, input, res, ctx)
	} else { // use cache if available
		// check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckResultCache(redisClient, input, ctx)
		// data not cached in Redis
		if (len(redisRes) == 0) {
			// acquire mutex on the key in Redis
			if ok, _ := redisClient.SetNX(ctx, input + "_mutex", 1, time.Duration(5) * time.Second).Result(); ok == true {
				// obtain analysis result and update Redis
				res = ObtainProfessor(input)
				RedisUpdateResultCache(redisClient, input, res, ctx)
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
				redisRes := RedisCheckResultCache(redisClient, input, ctx)
				// populate res with the retrieved data
				res = make([]string, 8)
				for ind, key := range profAttr {
					res[ind] = redisRes[key]
				}
			} 
		} else { // retrieve cached data from Redis
			// populate res with the retrieved data
			res = make([]string, 8)
			for ind, key := range profAttr {
				res[ind] = redisRes[key]
			}
		}
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		profAttr[0]: res[0], profAttr[1]: res[1], profAttr[2]: res[2], 
		profAttr[3]: res[3], profAttr[4]: res[4], profAttr[5]: res[5], 
		profAttr[6]: res[6], profAttr[7]: res[7], "queryHash": hash,
	})
}