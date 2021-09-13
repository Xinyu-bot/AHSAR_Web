package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	//"github.com/go-redis/redis/v8"
	//"context"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	r := gin.Default()

	// Simple MySQL connection
	var err error
	db, err = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/ahsar?parseTime=true&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect MySQL")
	}

	/*
		Periodically update MySQL database with newly added Professors on RMP
	*/
	go func() {
		for {
			// fetch max PID existed in MySQL database
			maxPID, err_max := ObtainMaxPID(db)
			if err_max != nil {
				log.Fatal("err_max:", err_max)
			}
			log.Println(maxPID)

			// initialize varibales for the loop
			var prof []Professor
			var err_prof error
			var nextPID string
			// fetch and update next professor info
			maxPID += 1
			nextPID = strconv.Itoa(maxPID)
			prof, err_prof = ObtainProfessorByPID(nextPID, "false", db)

			// while there is something new on RMP, continue to fetch more!
			for len(prof) == 1 && err_prof == nil {
				// fetch and update next professor info
				maxPID += 1
				nextPID = strconv.Itoa(maxPID)
				prof, err_prof = ObtainProfessorByPID(nextPID, "false", db)
			}

			// now that all newly added professors have been fetched and added to MySQL database,
			// compute next timestamp and wait
			now := time.Now()
			next := now.Add(time.Hour * 6)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()

	/*
		MiddleWare
	*/
	r.Use(CORSMiddleware())

	/*
		APIs avail for frontend
	*/
	// default GET for main page
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Welcome to AHSAR web page! "}) }) // return default response to frontend
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
