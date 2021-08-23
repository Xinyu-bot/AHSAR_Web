package main

import (
	"hash/crc32"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var profAttr = [7]string{"professor_name", "quality_score", "difficulty_score", "sentiment_score_discrete", "sentiment_score_continuous", "would_take_again", "pid"}

func main() {
	// 	Backend Server CORS setting
	r := gin.Default()

	// 	Simple Redis Connection
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		Password: "",  
		DB: 0,
	})
	defer redisClient.Close()

	/*
		MiddleWare
	*/
	r.Use(CORSMiddleware())

	/*
		APIs avail for frontend
	*/
	// 	default GET for main page
	r.GET("/", func(c *gin.Context) {
		// 	return default response to frontend
		c.JSON(200, gin.H{
			"message": "Welcome to AHSAR web page! ",
		})
	})

	// 	get professor sentiment analysis score by professor ID
	r.GET("/get_prof_by_id", func(c *gin.Context) {
		// 	extract input from request query
		input := c.Query("input")
		forceFetch := c.Query("noCache")
		// 	initialize res for storing result from NLP Server
		var res []string
		fmt.Println("[LOG] input:", input + ", forceFecth:", forceFetch)

		// if frontend enforces a no-cache-query, fetch latest data from RMP website
		if forceFetch == "true" {
			// 	obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
			res = ObtainProfessor(input)
			// 	update Redis with the new data
			if check := RedisCheckCache(redisClient, input); len(check) == 0 {
				RedisUpdateCache(redisClient, input, res)
			}
		} else { // use cache if available
			// 	check redis (as a fast RAM-based cache) if the input query exists
			redisRes := RedisCheckCache(redisClient, input)
			// 	data not cached in Redis
			if (len(redisRes) == 0) {
				// 	obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
				res = ObtainProfessor(input)
				// 	update Redis with the new data
				if check := RedisCheckCache(redisClient, input); len(check) == 0 {
					RedisUpdateCache(redisClient, input, res)
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
	})

	// 	serve on port 8080
	r.Run() 
}

// 	CORS for all origins
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Content-Type", "application/json")
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
        } else {
            c.Next()
        }
    }
}

//	fastHash function
func fastHash(ipAddr string, input string) string {
	hash := int(crc32.ChecksumIEEE([]byte(ipAddr + input)))
	res := strconv.Itoa(hash)
	return res
}