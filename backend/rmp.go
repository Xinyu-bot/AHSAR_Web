package main

import (
	"fmt"
	"bytes"
	"os/exec"
	"log"
	"strings"
)

func ObtainProfessor(input string) ([]string) {
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

	// if found
	outStr = strings.Trim(outStr, "\n")
	ret := strings.Split(outStr, " ")
	return ret
}