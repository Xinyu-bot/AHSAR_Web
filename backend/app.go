package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var profAttr = [6]string{"first_name", "last_name", "quality_score", "difficulty_score", "sentiment_score_discrete", "sentiment_score_continuous"}

func main() {
	r := gin.Default()

	// Simple Redis Connection
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		Password: "",  
		DB: 0,
	})
	defer redisClient.Close()

	/*
		APIs avail for frontend
	*/
	// default GET for main page
	r.GET("/", func(c *gin.Context) {
		// return default response to frontend
		c.JSON(200, gin.H{
			"message": "Welcome to AHSAR web page! ",
		})
	})

	// get professor sentiment analysis score by professor ID
	r.GET("/get_prof_by_id", func(c *gin.Context) {
		// extract input from request query
		input := c.Query("input")
		// initialize res for storing result from NLP Server
		var res []string

		// check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckCache(redisClient, input)
		// data not cached in Redis
		if (len(redisRes) == 0) {
			// obtain data of the professor from RMP website and analyze by ../pysrc/NLP_server.py
			res = ObtainProfessor(input)
			// update Redis with the new data
			RedisUpdateCache(redisClient, input, res)
		} else { // retrieve cached data from Redis
			// populate res with the retrieved data
			res = make([]string, 6)
			for ind, key := range profAttr {
				res[ind] = redisRes[key]
			}
		}

		// return response to frontend
		c.JSON(200, gin.H{
			profAttr[0]: res[0], profAttr[1]: res[1], profAttr[2]: res[2], 
			profAttr[3]: res[3], profAttr[4]: res[4], profAttr[5]: res[5],
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}