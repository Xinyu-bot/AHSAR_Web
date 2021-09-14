package main

import (
	"database/sql"
	"hash/crc32"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// port
var NLP_server_port = "localhost:5005"

// Abstruct Professor type
type Professor struct {
	PID                        string // `json: "pid" form: "pid`
	Prof_name                  string // `json: "prof_name" form:"prof_name`
	Quality_score              string // `json: "quality_score" form: "quality_score"`
	Difficulty_score           string // `json: "difficulty_score" form: "difficulty_score"`
	Would_take_again           string // `json: "would_take_again" form: "would_take_again"`
	School                     string // `json: "school" form: "school"`
	Department                 string // `json: "department" form: "department"`
	Sentiment_score_discrete   string // `json: "sentiment_score_discrete" form:"sentiment_score_discrete"`
	Sentiment_score_continuous string // `json: "sentiment_score_continuous" form:"sentiment_score_continuous"`
	Keywords                   string // `json: "keywords" form:"keywords"`
	Create_time                string // `json: "create_time" form: "create_time"`
	Update_time                string // `json: "upate_time" form: "upate_time"`
}

func PeriodicUpdate() {
	// fetch max PID existed in MySQL database
	maxPID, err_max := ObtainMaxPID(db)
	if err_max != nil {
		log.Fatal("err_max:", err_max)
	}
	/*
		check for newly added professors starting from the maxPID retrieved at the deployment of server,
		using the maxPID by the end of last periodic update, instead of dynamically fetching the maxPID from MySQL database...
		maybe some newly added professors have been inserted into the database because users have searched them by ID,
		but consider a case if we use dynamically fetched maxPID for each periodic update:
			1. assume PID #10 is the maxPID by the end of last periodic update, and #11 and after are all non-exist at that moment
			2. RMP adds #11 and #12 professors on their website
			3. user searches about #12 professor, so data about #12 professor is fetched from RMP, analyzed by NLP Server, and inserted into MySQL database
			4. the dynamically fetched maxPID is now #12, and next periodic update will start from #12 + 1 = #13
			5. notice that #11 is skipped and won't be able to be searched by users in SDP or by name, until someone searches #11 by ID
	*/
	for {
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
}

// Obtain maximum PID in the Database
func ObtainMaxPID(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT MAX(pid) FROM prof")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ret int
	for rows.Next() {
		err = rows.Scan(&ret)
		if err != nil {
			log.Fatal(err)
		}
	}

	return ret, err
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
	if err_t != nil {
		log.Fatal(err_t)
	}
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
		return ret_t, err_t
	}

	// if noCache is true, then MUST fetch latest data from RMP website, analyze, and update MySQL
	if noCache == "true" || flag {
		// Simple TCP Connection to NLP Server
		conn, errTCP := net.Dial("tcp", NLP_server_port)
		if errTCP != nil {
			log.Fatal("errTCP:", errTCP)
		}
		defer conn.Close()

		// send input to NLP Server
		// 0 is defined internally here meaning query parameter is PID and wants sentiment analysis result
		conn.Write([]byte("0" + input))

		// read the analyzed result from NLP Server
		buf := make([]byte, 4096)
		n, errRead := conn.Read(buf)
		if errRead != nil {
			log.Fatal("errRead:", errRead)
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
			stmt, errQuery := db.Prepare("INSERT INTO prof (pid, prof_name, quality_score, difficulty_score, would_take_again, school, department, sentiment_score_continuous, sentiment_score_discrete, keywords) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE quality_score = ?, difficulty_score = ?, would_take_again = ?, sentiment_score_continuous = ?, sentiment_score_discrete = ?, keywords = ?, update_time = current_timestamp")
			if errQuery != nil {
				log.Fatal("errQuery:", errQuery)
			}
			// execute the insert query
			_, errExec := stmt.Exec(ret[7], ret[0], ret[1], ret[2], ret[6], ret[8], ret[9], ret[4], ret[3], ret[5], // INSERT
				ret[1], ret[2], ret[6], ret[4], ret[3], ret[5], // UPDATE
			)
			if errExec != nil {
				log.Fatal("errExec:", errExec)
			}
		} else { // if professor does not exist, return EMPTY []Professor{} list and nil error
			ret := []Professor{}
			return ret, nil
		}
	}

	// Obtain Professor by PID
	rows, err := db.Query("SELECT * FROM prof WHERE pid = ?", input)
	if err != nil {
		log.Fatal(err)
	}
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
	rows, err := db.Query("SELECT pid, prof_name, quality_score, difficulty_score, would_take_again, school, department, sentiment_score_continuous, sentiment_score_discrete, keywords, create_time, update_time FROM prof WHERE MATCH(prof_name) AGAINST(?) LIMIT 20", name)
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		log.Fatal(err)
	}
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
	rows, err := db.Query("SELECT DISTINCT school FROM prof WHERE school LIKE ?", initial+"%")
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		log.Fatal(err)
	}
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
