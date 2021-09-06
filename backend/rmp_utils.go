package main

import (
	"net"
	"log"
	"strings"
	"database/sql"
	"hash/crc32"
	"strconv"
)

var profAttr = [8]string{"professor_name", "quality_score", "difficulty_score", "sentiment_score_discrete", "sentiment_score_continuous", "keywords", "would_take_again", "pid"}
var profPIDAttr = [4]string{"professor_name", "department", "school", "pid"}
var NLP_server_port = "localhost:5005"
// Abstruct Professor type
type Professor struct {
	PID 				string 	`json: "pid" form: "pid`
	Prof_name 			string 	`json: "prof_name" form:"prof_name`
	Quality_score		string 	`json: "quality_score" form: "quality_score"`
	Difficulty_score 	string 	`json: "difficulty_score" form: "difficulty_score"`
	Would_take_again 	string 	`json: "would_take_again" form: "would_take_again"`
	School 				string 	`json: "school" form: "school"`	
	Department 			string 	`json: "department" form: "department"`
}

// Obtain Professor result by PID
func ObtainProfessor(input string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", NLP_server_port)
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	// 0 is defined internally here meaning query parameter is PID and wants sentiment analysis result
	conn.Write([]byte("0" + input))

	// read result from NLP Server
	buf := make([]byte, 4096)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), " ")
	ret[0] = strings.Replace(ret[0], "?", " ", -1)
	ret[5] = strings.Replace(ret[5], "`", " ", -1)
	return ret
}

// Obtain professor list by name
func ObtainProfessorByName(name string, db *sql.DB) ([]Professor, error) {
	rows, err := db.Query("SELECT pid, prof_name, quality_score, difficulty_score, would_take_again, school, department FROM prof WHERE MATCH(prof_name) AGAINST(?) LIMIT 10", name)
	defer rows.Close()

	ret := []Professor{}
	for rows.Next() {
		professor := Professor{}
        err = rows.Scan(&professor.PID, &professor.Prof_name, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again, &professor.School, &professor.Department)
		if err != nil {
            log.Fatal(err)
        }
		ret = append(ret, professor)
	}

	return ret, err
}

// Obtain departments list by school name
func ObtainDepartmentList(school string, db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT department FROM prof WHERE school = ?", school)
	defer rows.Close()

	ret := []string{}
	for rows.Next() {
		var department string
        err = rows.Scan(&department)
		if err != nil {
            log.Fatal(err)
        }
		ret = append(ret, department)
	}

	return ret, err
}

// Obtain School list by letter initial
func ObtainSchoolList(initial string, db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT school FROM prof WHERE school LIKE ?", initial + "%")
	defer rows.Close()

	ret := []string{}
	for rows.Next() {
		var school string
        err = rows.Scan(&school)
		if err != nil {
            log.Fatal(err)
        }
		ret = append(ret, school)
	}

	return ret, err
}

// Obtain professors list by school and department name
func ObtainProfessorList(school, department string, db *sql.DB) ([]Professor, error) {
	rows, err := db.Query("SELECT pid, prof_name, quality_score, difficulty_score, would_take_again, school, department FROM prof WHERE school = ? and department = ? ORDER BY quality_score DESC, difficulty_score ASC", school, department)
	defer rows.Close()

	ret := []Professor{}
	for rows.Next() {
		professor := Professor{}
        err = rows.Scan(&professor.PID, &professor.Prof_name, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again, &professor.School, &professor.Department)
		if err != nil {
            log.Fatal(err)
        }
		ret = append(ret, professor)
	}

	return ret, err
}

// fastHash function
func fastHash(ipAddr string, input string) string {
	hash := int(crc32.ChecksumIEEE([]byte(ipAddr + input)))
	res := strconv.Itoa(hash)
	return res
}