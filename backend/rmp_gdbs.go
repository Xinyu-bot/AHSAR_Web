package main

import (
	"log"
	//"time"
	"github.com/gin-gonic/gin"
)

func GetDepartmentsBySchool(c *gin.Context) {
	// extract school name from request query
	school := c.Query("school")

	// initialize hasResult storing boolean value whether there is actual result or not
	var hasResult string
	log.Println("query input: school name = {" + school + "}")

	/*
		fetch from MySQL directly using indexed qeury
		response time in 5ms, which is fast enough
		so as for now, no Redis (as cache database) is involved
	*/
	ret, err := ObtainDepartmentList(school, db)
	if err != nil || len(ret) == 0 {
		hasResult = "false"
	} else {
		hasResult = "true"
	}

	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), school)
	c.JSON(200, gin.H{
		"queryHash":      hash,
		"hasResult":      hasResult,
		"departmentList": ret,
	})
}
