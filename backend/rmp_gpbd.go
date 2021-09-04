package main

import (
	"log"
	// "time"
	"github.com/gin-gonic/gin"
)

func GetProfsByDepartment(c *gin.Context) {
	// extract school name from request query
	school := c.Query("school")
	department := c.Query("department")
	
	// initialize res for storing result from NLP Server
	// var res []string
	log.Println("query input: school = {" + school + "}, department = {" + department + "}")

	// TODO


	// return response to frontend with unique hash for each query
	hash := fastHash(c.ClientIP(), school + department)
	c.JSON(400, gin.H{
		"queryHash": hash,
	})
}