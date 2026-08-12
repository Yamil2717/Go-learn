package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	countLines := false
	if len(os.Args) > 1 && os.Args[1] == "-l" {
		countLines = true
	}

	scanner := bufio.NewScanner(os.Stdin)
	wordCount := 0
	lineCount := 0

	for scanner.Scan() {
		text := scanner.Text()

		if strings.TrimSpace(text) == "exit" {
			break
		}

		lineCount++

		words := strings.Fields(text)
		wordCount += len(words)
	}

	if countLines {
		fmt.Printf("Total lines: %d\n", lineCount)
	} else {
		fmt.Printf("Total words: %d\n", wordCount)
	}
}
