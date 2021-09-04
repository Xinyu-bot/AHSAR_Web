package main

import (
	"log"
	// "time"
	"github.com/gin-gonic/gin"
)

func GetSchoolsByInitial(c *gin.Context) {
	// extract school name from request query
	initial := c.Query("initial")
	
	// initialize res for storing result from NLP Server
	// var res []string
	log.Println("query input: initial = {" + initial + "}")

	// TODO


	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), initial)
	c.JSON(400, gin.H{
		"queryHash": hash,
	})
}