package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

func GetDepartmentsBySchool(c *gin.Context) {
	// extract school name from request query
	school := c.Query("school")
	
	// initialize res for storing result from NLP Server
	var res, ret []string
	var hasResult string
	log.Println("query input: school name = {" + school + "}")

	// check redis (as a fast RAM-based cache) if the input query exists
	redisRes := RedisCheckDepartmentList(redisClient, school, ctx)
	// data not cached in Redis
	if (len(redisRes) == 0) {
		// acquire mutex on the key in Redis
		if ok, _ := redisClient.SetNX(ctx, "school" + school + "_mutex", 1, time.Duration(5) * time.Second).Result(); ok == true {
			// obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
			res = ObtainDepartments(school)
			if res[0] == "-1" {
				hasResult = "false"
				ret = append(ret, "DNE")
			} else {
				hasResult = "true"
				ret = res
			}
			RedisUpdateDepartmentList(redisClient, "school" + school, ret, ctx)
			// release mutex
			redisClient.Del(ctx, "school" + school + "_mutex")
		} else {
			// fail to acquire, meaning other goroutine is holding the mutex... then sleep for 50ms
			time.Sleep(50)
			// check mutex every 50ms
			for ;; {
				checkMutex, _ := redisClient.Exists(ctx, "school" + school + "_mutex").Result()
				if checkMutex == 1 {
					time.Sleep(50)
				} else { // if mutex has been released, 
					// meaning that Redis should have cached data, then break and get the cached data
					break
				}
			}
			// obtain cached data from Redis
			redisRes = RedisCheckDepartmentList(redisClient, "school" + school, ctx)
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
	hash := fastHash(c.ClientIP(), school)
	c.JSON(200, gin.H{
		"queryHash": hash,
		"hasResult": hasResult,
		"dList": ret,
	})
}