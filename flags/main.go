// DON'T CHANGE THIS FILE
package main

import (
	"flags/fl"
	"fmt"
)

func main() {
	lines := fl.Bool("-l", false, "Flag to determine if the program should count words or lines")
	fmt.Println("Before parsing:", *lines) // Prints default value --> false
	fl.Parse()
	fmt.Println("After parsing:", *lines) // Prints parsed value

	if *lines {
		fmt.Println("The program should count lines")
	} else {
		fmt.Println("The program should count words")
	}
}
