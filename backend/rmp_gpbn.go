package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetPidByName(c *gin.Context) {
	// extract input from request query
	input := c.Query("input")
	// initialize hasResult storing boolean value whether there is actual result or not
	var hasResult string
	fmt.Println("[LOG] query input: " + input)

	/*
		fetch from MySQL directly using indexed qeury
		response time in 5ms, which is fast enough
		so as for now, no Redis (as cache database) is involved
	*/
	ret, err := ObtainProfessorByName(input, db)
	if err != nil || len(ret) == 0 {
		hasResult = "false"
	} else {
		hasResult = "true"
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), input)
	c.JSON(200, gin.H{
		"queryHash": hash,
		"hasResult": hasResult,
		"ret":       ret,
	})
}
