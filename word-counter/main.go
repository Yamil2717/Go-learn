package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func count(input io.Reader, countLines bool) int {
	scanner := bufio.NewScanner(input)
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
		return lineCount
	}
	return wordCount
}

func main() {
	countLines := false
	if len(os.Args) > 1 && os.Args[1] == "-l" {
		countLines = true
	}

	result := count(os.Stdin, countLines)

	if countLines {
		fmt.Printf("Total lines: %d\n", result)
	} else {
		fmt.Printf("Total words: %d\n", result)
	}
}
