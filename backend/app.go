package main

import (
	"log"
	//"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	//"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	r := gin.Default()

	// Simple MySQL connection
	var err error
	db, err = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/ahsar?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil{
		log.Fatal("Failed to connect MySQL")
	}

	/*
		MiddleWare
	*/
	r.Use(CORSMiddleware())

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
	r.GET("/get_prof_by_id", GetProfByID)
	// get professor's PID from RMP website by professor name
	r.GET("/get_pid_by_name", GetPidByName)
	// TOOD: get schools list by initial
	r.GET("/get_schools_by_initial", GetSchoolsByInitial)
	// get departments list by school name
	r.GET("/get_departments_by_school", GetDepartmentsBySchool)
	// get professors list by school and department names
	r.GET("/get_prof_by_department", GetProfByDepartment)

	// serve on port 8080
	r.Run(":8080") 
}

// CORS for all origins
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