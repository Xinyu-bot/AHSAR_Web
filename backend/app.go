package main

import (
	"fmt"
	//"net/http"
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to AHSAR web page! ",
		})
	})

	// get professor sentiment analysis score by professor ID
	r.GET("/get_prof_by_id", func(c *gin.Context) {
		input := c.Query("input")
		fmt.Println("Client IP", c.ClientIP)
		var res []string

		// check redis (as a fast RAM-based cache) if the input query exists
		redisRes := RedisCheckCache(redisClient, input)
		if (len(redisRes) == 0) {
			// data not cached in Redis
			fmt.Println(input, "is not in Redis")
			/*
			obtain data of the professor from RMP website and analyze:
				all done in Python code under ./pysrc/application.py
				res is in form of 
				[]string{name, quality_score, difficulty_score, sentiment_score_discrete, sentiment_score_continuous}
			*/
			res = ObtainProfessor(input)
			fmt.Println("res from ObtainProfessor:", res)
			// update Redis with the new data
			RedisUpdateCache(redisClient, input, res)
		} else {
			// retrieve cached data from Redis
			fmt.Println(input, "is in Redis")
			// populate res with the retrieved data
			res = make([]string, 6)
			for ind, key := range profAttr {
				res[ind] = redisRes[key]
			}
		}

		fmt.Println("returned message:", res)
		c.JSON(200, gin.H{
			"messgae": res,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}