package main

import (
	"net"
	"log"
	"strings"
)

func ObtainProfessor(input string) ([]string) {
	// 	Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", "localhost:5005")
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// 	send input to NLP Server
	conn.Write([]byte(input))

	// 	read result from NLP Server
	buf := make([]byte, 2048)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// 	format the read result
	ret := strings.Split(string(buf[:n]), " ")
	return ret
}