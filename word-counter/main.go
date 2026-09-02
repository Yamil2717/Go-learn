package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func capture(input io.Reader) []string {
	scanner := bufio.NewScanner(input)
	var lines []string

	for scanner.Scan() {
		text := scanner.Text()

		if strings.TrimSpace(text) == "exit" {
			break
		}

		lines = append(lines, text)
	}

	return lines
}

func count(lines []string, countLines bool) int {
	if countLines {
		return len(lines)
	}

	total := 0
	for _, line := range lines {
		total += len(strings.Fields(line))
	}

	return total
}

func main() {
	countLines := len(os.Args) > 1 && os.Args[1] == "-l"

	lines := capture(os.Stdin)
	result := count(lines, countLines)

	if countLines {
		fmt.Printf("Total lines: %d\n", result)
	} else {
		fmt.Printf("Total words: %d\n", result)
	}
}
