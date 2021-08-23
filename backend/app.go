package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

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
	r.GET("/get_prof_by_id", GetProfByID)
	//	get professor's PID from RMP website by professor name
	r.GET("/get_pid_by_name", GetPidByName)

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