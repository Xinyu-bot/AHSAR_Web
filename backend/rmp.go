package main

import (
	"fmt"
	//"bytes"
	"net"
	//"os/exec"
	"log"
	"strings"
	//"strconv"
	//"github.com/Shopify/sarama"
)

func ObtainProfessor(/*producer sarama.AsyncProducer, */seq uint64, input string) ([]string) {

	/*
	producer.Input() <- &sarama.ProducerMessage{Topic: "NLP", Key: nil, Value: sarama.StringEncoder(input + "," + strconv.FormatUint(seq, 10))}
	// wait response
	select {
		case msg := <-producer.Successes():
			log.Printf("Produced message successes: [%s]\n",msg.Value)
		case err := <-producer.Errors():
			log.Println("Produced message failure: ", err)
	}
	*/

	// Simple TCP Connection
	conn, errTCP := net.Dial("tcp", "localhost:5000")
	if errTCP != nil {
		log.Fatalf("errTCP:", errTCP)
	}
	defer conn.Close()

	fmt.Println(conn)
	conn.Write([]byte(input))

	buf := make([]byte, 1024)
	n, errRead := conn.Read(buf)
	if errRead != nil {
		log.Fatalf("errRead:", errRead)
	}
	fmt.Println(string(buf[:n]))

	/*
	// Old way of running application.py script -> one python process per GET request
	args := []string{"../pysrc/application.py", input}
	cmd := exec.Command("python3", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err, errStr)
    }

	// if professor is not found or no comments
	fmt.Println("outStr in rmp.go/ObtainProfessor:", outStr)
	if len(outStr) == 0 {
		return []string{"-1", "-1", "-1", "-1", "-1", "-1"}
	}
	*/

	// if found
	// outStr = strings.Trim(outStr, "\n")
	ret := strings.Split(string(buf[:n]), " ")
	return ret
}