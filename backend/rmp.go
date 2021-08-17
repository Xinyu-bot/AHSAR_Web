package main

import (
	"fmt"
	"bytes"
	"os/exec"
	"log"
	"strings"
)

func ObtainProfessor(input string) ([]string) {
	var ret []string

	args := []string{"../pysrc/application.py", input}
	cmd := exec.Command("python3", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err, errStr)
    }

	// if professor is not found 
	fmt.Println("outStr in rmp.go/ObtainProfessor:", outStr)
	if len(outStr) == 0 {
		return []string{"-1", "-1", "-1", "-1", "-1", "-1"}
	}

	outStr = strings.Trim(outStr, "\n")
	ret = strings.Split(outStr, " ")
	return ret
}