package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/secure"
	//"github.com/go-redis/redis/v8"
)

var db *sql.DB

func main() {
	// set gin Framework to release mode
	gin.SetMode(gin.ReleaseMode)
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
	go PeriodicUpdate()

	/*
		MiddleWare
	*/
	r.Use(TlsHandler())
	r.Use(CORSMiddleware())

	/*
		default GET for main page
	*/
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Welcome to AHSAR web page! "}) }) // return default response to frontend

	/*
		Static File
	*/
	r.StaticFile("/favicon.ico", "./src/static/favicon/favicon.ico")
	r.StaticFile("/logo.png", "./src/static/logo/logo.png")
	r.StaticFile("/logo.svg", "./src/static/logo/logo.svg")

	/*
		Actual Service
	*/
	r.GET("/get_prof_by_id", GetProfByID)                       // get professor sentiment analysis score by professor ID
	r.GET("/get_pid_by_name", GetPidByName)                     // get professor's PID from RMP website by professor name
	r.GET("/get_schools_by_initial", GetSchoolsByInitial)       // get schools list by initial
	r.GET("/get_departments_by_school", GetDepartmentsBySchool) // get departments list by school name
	r.GET("/get_prof_by_department", GetProfByDepartment)       // get professors list by school and department names

	// serve on port 8080
	r.RunTLS(":8080", "./ahsar.pem", "./ahsar.key")
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

// TSL handler #https://www.jianshu.com/p/01057d2c37e4
func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			log.Fatal("TLS Handler Failure...")
		}
		c.Next()
	}
}
