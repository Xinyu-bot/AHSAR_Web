package main

import (
	"fmt"
	//"log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	//"github.com/Shopify/sarama"
	"github.com/segmentio/fasthash/fnv1a"
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
	// Simple Kafka Connection
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("unable to create kafka client: %q", err)
	}
	defer kafkaClient.Close()
	producer, errk := sarama.NewAsyncProducerFromClient(kafkaClient)
	if errk != nil {
		log.Fatalf("unable to create kafka producer: %q", errk)
	}
	defer producer.Close()
	*/

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
			seq := fnv1a.HashString64(input)
			res = ObtainProfessor(/*producer, */seq, input)
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
			profAttr[0]: res[0],
			profAttr[1]: res[1],
			profAttr[2]: res[2],
			profAttr[3]: res[3],
			profAttr[4]: res[4],
			profAttr[5]: res[5],
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}