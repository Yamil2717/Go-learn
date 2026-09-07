package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	cmd := exec.Command("ls", "-lah")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() fallo con %s\n", err)
	}
	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	fmt.Printf("\nSalida:\n%s\nError:\n%s\n", outStr, errStr)
}
