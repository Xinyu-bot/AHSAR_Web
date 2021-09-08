package main

import (
	"net"
	"log"
	"strings"
	"database/sql"
	"hash/crc32"
	"strconv"
)

// port
var NLP_server_port = "localhost:5005"
// Abstruct Professor type
type Professor struct {
	PID 						string 	`json: "pid" form: "pid`
	Prof_name 					string 	`json: "prof_name" form:"prof_name`
	Quality_score				string 	`json: "quality_score" form: "quality_score"`
	Difficulty_score 			string 	`json: "difficulty_score" form: "difficulty_score"`
	Would_take_again 			string 	`json: "would_take_again" form: "would_take_again"`
	School 						string 	`json: "school" form: "school"`	
	Department 					string 	`json: "department" form: "department"`
	Sentiment_score_discrete 	string 	`json: "sentiment_score_discrete" form:"sentiment_score_discrete"`
	Sentiment_score_continuous 	string 	`json: "sentiment_score_continuous" form:"sentiment_score_continuous"`
	Keywords					string 	`json: "keywords" form:"keywords"`
	Create_time 				string 	`json: "create_time" form: "create_time"`
	Update_time					string 	`json: "upate_time" form: "upate_time"`
}

// Obtain Professor result by PID
func ObtainProfessorByPID(input, noCache string, db *sql.DB) ([]Professor, error) {
	/*
		brutal way of checking MySQL has sentiment analysis result or not
		if not, try to fetch data from RMP
	*/
	flag := false
	
	// Obtain Professor by PID
	rows_t, err_t := db.Query("SELECT * FROM prof WHERE pid = ?", input)
	defer rows_t.Close()

	ret_t := []Professor{}
	for rows_t.Next() {
		professor := Professor{}
        err_t = rows_t.Scan(
			&professor.PID, &professor.Prof_name, &professor.School, &professor.Department, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again,  
			&professor.Sentiment_score_continuous, &professor.Sentiment_score_discrete, &professor.Keywords, &professor.Create_time, &professor.Update_time, 
		)
		if err_t != nil {
            log.Fatal(err_t)
        }
		ret_t = append(ret_t, professor)
	}
	// if professor does not exist in the table or his sentiment scores have not yet been updated
	if len(ret_t) == 0 || ret_t[0].Sentiment_score_discrete == "-2" {
		flag = true
	} 
	// if all good... return the result
	if len(ret_t) == 1 && ret_t[0].Sentiment_score_discrete != "-2" && noCache != "true" {
		log.Println(ret_t[0])
		return ret_t, err_t
	}

	// if noCache is true, then MUST fetch latest data from RMP website, analyze, and update MySQL
	if noCache == "true" || flag == true {
		// Simple TCP Connection to NLP Server
		conn, errTCP := net.Dial("tcp", NLP_server_port)
		if errTCP != nil {
			log.Fatalf("errTCP:", errTCP)
		}
		defer conn.Close()

		// send input to NLP Server
		// 0 is defined internally here meaning query parameter is PID and wants sentiment analysis result
		conn.Write([]byte("0" + input))

		// read the analyzed result from NLP Server
		buf := make([]byte, 4096)
		n, errRead := conn.Read(buf)
		if errRead != nil {
			log.Fatalf("errRead:", errRead)
		}

		// format the result
		ret := strings.Split(string(buf[:n]), " ")
		ret[0] = strings.Replace(ret[0], "?", " ", -1)
		ret[5] = strings.Replace(ret[5], "`", " ", -1)
		ret[8] = strings.Replace(ret[8], "`", " ", -1)
		ret[9] = strings.Replace(ret[9], "`", " ", -1)

		// update the MySQL database only if the professor does exist, i.e. PID is not -1 (Check `../pysrc/NLP_server.py` for detail about the returned values)
		if ret[0] != "-1" {
			/*
				prepare the insert query first:
					1. if professor is already in the table, update scores and keywords
					2. if not, insert as a new row
			*/
			stmt, errQuery := db.Prepare("INSERT INTO prof (pid, prof_name, quality_score, difficulty_score, would_take_again, school, department, sentiment_score_continuous, sentiment_score_discrete, keywords) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE quality_score = ?, difficulty_score = ?, would_take_again = ?, sentiment_score_continuous = ?, sentiment_score_discrete = ?, keywords = ?")
			if errQuery != nil {
				log.Fatalf("errQuery:", errQuery)
			}
			// execute the insert query
			_, errExec := stmt.Exec(ret[7], ret[0], ret[1], ret[2], ret[6], ret[8], ret[9], ret[4], ret[3], ret[5], // INSERT
									ret[1], ret[2], ret[6], ret[4], ret[3], ret[5], // UPDATE
									) 
			if errExec != nil {
				log.Fatalf("errExec:", errExec)
			}
		} else { // if professor does not exist, return EMPTY []Professor{} list and nil error
			ret := []Professor{}
			return ret, nil
		}
	}

	// Obtain Professor by PID
	rows, err := db.Query("SELECT * FROM prof WHERE pid = ?", input)
	defer rows.Close()

	ret := []Professor{}
	for rows.Next() {
		professor := Professor{}
        err = rows.Scan(
			&professor.PID, &professor.Prof_name, &professor.School, &professor.Department, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again, 
			&professor.Sentiment_score_continuous, &professor.Sentiment_score_discrete, &professor.Keywords, &professor.Create_time, &professor.Update_time, 
		)
		if err != nil {
            log.Fatal(err)
        }
		ret = append(ret, professor)
	}

	return ret, err
}

// Obtain professor list by name
func ObtainProfessorByName(name string, db *sql.DB) ([]Professor, error) {
	rows, err := db.Query("SELECT pid, prof_name, quality_score, difficulty_score, would_take_again, school, department, sentiment_score_continuous, sentiment_score_discrete, keywords, create_time, update_time FROM prof WHERE MATCH(prof_name) AGAINST(?) LIMIT 10", name)
	defer rows.Close()

	ret := []Professor{}
	for rows.Next() {
		professor := Professor{}
        err = rows.Scan(
			&professor.PID, &professor.Prof_name, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again, &professor.School, &professor.Department, 
			&professor.Sentiment_score_continuous, &professor.Sentiment_score_discrete, &professor.Keywords, &professor.Create_time, &professor.Update_time, 
		)
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
	rows, err := db.Query("SELECT pid, prof_name, quality_score, difficulty_score, would_take_again, school, department, sentiment_score_continuous, sentiment_score_discrete, keywords, create_time, update_time FROM prof WHERE school = ? and department = ? ORDER BY quality_score DESC, difficulty_score ASC", school, department)
	defer rows.Close()

	ret := []Professor{}
	for rows.Next() {
		professor := Professor{}
        err = rows.Scan(
			&professor.PID, &professor.Prof_name, &professor.Quality_score, &professor.Difficulty_score, &professor.Would_take_again, &professor.School, &professor.Department, 
			&professor.Sentiment_score_continuous, &professor.Sentiment_score_discrete, &professor.Keywords, &professor.Create_time, &professor.Update_time, 
		)
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