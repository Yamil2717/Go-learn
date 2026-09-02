package main

import (
	"bufio"
	"flag"
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

func count(lines []string, countLines bool, countBytes bool) int {
	if countLines {
		return len(lines)
	}

	if countBytes {
		return len(strings.Join(lines, "\n"))
	}

	total := 0
	for _, line := range lines {
		total += len(strings.Fields(line))
	}

	return total
}

func main() {
	countLines := flag.Bool("l", false, "Flag to determine if the program should count lines")
	countBytes := flag.Bool("b", false, "Flag to determine if the program should count bytes")
	flag.Parse()

	lines := capture(os.Stdin)
	result := count(lines, *countLines, *countBytes)

	if *countLines {
		fmt.Printf("Total lines: %d\n", result)
	} else if *countBytes {
		fmt.Printf("Total bytes: %d\n", result)
	} else {
		fmt.Printf("Total words: %d\n", result)
	}
}
