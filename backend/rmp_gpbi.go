package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetProfByID(c *gin.Context) {
	// extract input from request query
	input := c.Query("input")
	noCache := c.Query("noCache")
	// initialize hasResult storing boolean value whether there is actual result or not
	var hasResult string
	// initialize res for storing result
	var res Professor
	fmt.Println("[LOG] query input: " + input + ", noCache: " + noCache)

	/*
		fetch from MySQL directly using indexed qeury
		response time in 5ms, which is fast enough
		so as for now, no Redis (as cache database) is involved
	*/
	ret, err := ObtainProfessorByPID(input, noCache, db)
	if err != nil || len(ret) == 0 {
		hasResult = "false"
		res = Professor{}
	} else {
		hasResult = "true"
		res = ret[0]
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		"queryHash": hash,
		"hasResult": hasResult,
		"professor": res,
	})
}
