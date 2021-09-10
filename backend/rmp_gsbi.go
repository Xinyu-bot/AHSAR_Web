package main

import (
	"log"
	//"time"
	"github.com/gin-gonic/gin"
)

func GetSchoolsByInitial(c *gin.Context) {
	// extract school name from request query
	initial := c.Query("initial")

	// initialize hasResult storing boolean value whether there is actual result or not
	var hasResult string
	log.Println("query input: initial = {" + initial + "}")

	/*
		fetch from MySQL directly using indexed qeury
		response time in 5ms, which is fast enough
		so as for now, no Redis (as cache database) is involved
	*/
	ret, err := ObtainSchoolList(initial, db)
	if err != nil || len(ret) == 0 {
		hasResult = "false"
	} else {
		hasResult = "true"
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), initial)
	c.JSON(200, gin.H{
		"queryHash":  hash,
		"schoolList": ret,
		"hasResult":  hasResult,
	})
}
