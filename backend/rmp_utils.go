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

// Obtain Professor result by PID
func ObtainProfessor(input string) ([]string) {
	// 	Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", "localhost:5005")
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// 	send input to NLP Server
	conn.Write([]byte("0"+input))

	// 	read result from NLP Server
	buf := make([]byte, 2048)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// 	format the read result
	ret := strings.Split(string(buf[:n]), " ")
	ret[0] = strings.Replace(ret[0], "?", " ", -1)
	return ret
}

// 	Obatin professor's PID by name
func ObtainPID(input string) ([]string) {
	// 	Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", "localhost:5005")
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// 	send input to NLP Server
	conn.Write([]byte("1"+input))

	// 	read result from NLP Server
	buf := make([]byte, 4096)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// 	format the read result
	ret := strings.Split(string(buf[:n]), " ")
	return ret
}


//	fastHash function
func fastHash(ipAddr string, input string) string {
	hash := int(crc32.ChecksumIEEE([]byte(ipAddr + input)))
	res := strconv.Itoa(hash)
	return res
}