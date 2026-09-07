package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func main() {
	cmd := exec.Command("ls", "-lah")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", b.String())
}
