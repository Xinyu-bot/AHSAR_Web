package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

func GetSchoolsByInitial(c *gin.Context) {
	// extract school name from request query
	initial := c.Query("initial")
	
	// initialize res for storing result from NLP Server
	var res, ret []string
	var hasResult string
	log.Println("query input: initial = {" + initial + "}")

	// check redis (as a fast RAM-based cache) if the input query exists
	redisRes := RedisCheckDepartmentList(redisClient, "initial" + initial, ctx)
	// data not cached in Redis
	if (len(redisRes) == 0) {
		// acquire mutex on the key in Redis
		if ok, _ := redisClient.SetNX(ctx, "initial" + initial + "_mutex", 1, time.Duration(5) * time.Second).Result(); ok == true {
			// obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
			res = ObtainSchoolList(initial)
			if res[0] == "-1" {
				hasResult = "false"
				ret = append(ret, "DNE")
			} else {
				hasResult = "true"
				ret = res
			}
			// Reuse utils function for Department List here --> Maybe rename this total shyt
			RedisUpdateDepartmentList(redisClient, "initial" + initial, ret, ctx)
			// release mutex
			redisClient.Del(ctx, "initial" + initial + "_mutex")
		} else {
			// fail to acquire, meaning other goroutine is holding the mutex... then sleep for 50ms
			time.Sleep(50)
			// check mutex every 50ms
			for ;; {
				checkMutex, _ := redisClient.Exists(ctx, "initial" + initial + "_mutex").Result()
				if checkMutex == 1 {
					time.Sleep(50)
				} else { // if mutex has been released, 
					// meaning that Redis should have cached data, then break and get the cached data
					break
				}
			}
			// obtain cached data from Redis
			redisRes = RedisCheckDepartmentList(redisClient, "initial" + initial, ctx)
			// populate res with the retrieved data
			if res[0] == "-1" {
				hasResult = "false"
				ret = append(ret, "DNE")
			} else {
				hasResult = "true"
				ret = res
			}
		}
	} else { //  retrieve cached data from Redis
		// populate res with the retrieved data
		if redisRes["0"] == "-1" {
			hasResult = "false"
			ret = append(ret, "DNE")
		} else {
			hasResult = "true"
			for _, val := range redisRes {
				ret = append(ret, val)
			}
		}
	}


	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), initial)
	c.JSON(200, gin.H{
		"queryHash": hash,
		"pList": ret,
		"hasResult": hasResult,
	})
}