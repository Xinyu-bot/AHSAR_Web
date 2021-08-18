package main

import (
	"net"
	"log"
	"strings"
)

func ObtainProfessor(input string) ([]string) {
	// Simple TCP Connection to NLP Server
	conn, errTCP := net.Dial("tcp", "localhost:5000")
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	// send input to NLP Server
	conn.Write([]byte(input))

	buf := make([]byte, 1024)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}

	// read result from NLP Server
	ret := strings.Split(string(buf[:n]), " ")
	return ret
}