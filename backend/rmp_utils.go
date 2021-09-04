package main

import (
	"net"
	"log"
	"strings"
	"hash/crc32"
	"strconv"
)

var profAttr = [7]string{"professor_name", "quality_score", "difficulty_score", "sentiment_score_discrete", "sentiment_score_continuous", "would_take_again", "pid"}
var profPIDAttr = [4]string{"professor_name", "department", "school", "pid"}
var NLP_server_port = "localhost:5005"

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
	buf := make([]byte, 2048)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), " ")
	ret[0] = strings.Replace(ret[0], "?", " ", -1)
	return ret
}

// Obatin professor's PID by name
func ObtainPID(input string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", NLP_server_port)
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	// 1 is defined internally here meaning query parameter is Name and wants list of relevant professor info
	conn.Write([]byte("1" + input))

	// read result from NLP Server
	buf := make([]byte, 4096)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), " ")
	return ret
}

// Obtain departments list by school name
func ObtainDepartments(school string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", NLP_server_port)
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	// 2 is defined internally here meaning query parameter is school name and wants departments list
	conn.Write([]byte("2" + school))

	// read result from NLP Server
	buf := make([]byte, 32768)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), "$")
	return ret
}

// Obtain professors list by school and department name
func ObtainProfessorList(school, department string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", NLP_server_port)
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	// 3 is defined internally here meaning 
	// query parameter is school and department name and wants departments list
	conn.Write([]byte("3" + school + "$" + department))

	// read result from NLP Server
	buf := make([]byte, 32768)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), "$")
	return ret
}

// Obtain School list by letter initial
func ObtainSchoolList(initial string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", NLP_server_port)
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	// 4 is defined internally here meaning 
	// query parameter is letter initial and wants schools list
	conn.Write([]byte("4" + initial))

	// read result from NLP Server
	buf := make([]byte, 32768)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// format the read result
	ret := strings.Split(string(buf[:n]), "$")
	return ret
}

// fastHash function
func fastHash(ipAddr string, input string) string {
	hash := int(crc32.ChecksumIEEE([]byte(ipAddr + input)))
	res := strconv.Itoa(hash)
	return res
}