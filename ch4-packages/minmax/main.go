package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func minmax(min float64, max float64, values ...float64) []float64 {
	var result []float64

	for _, value := range values {
		if value >= min && value <= max {
			result = append(result, value)
		}
	}

	return result
}

func main() {
	fmt.Print("Ingresa el numero minimo: ")
	minInput := getInput()

	fmt.Print("Ingresa el numero maximo: ")
	maxInput := getInput()

	fmt.Print("Ingresa los numeros separados por espacio: ")
	valuesInput := getInput()

	min, _ := strconv.ParseFloat(minInput, 64)
	max, _ := strconv.ParseFloat(maxInput, 64)

	valueStrings := strings.Fields(valuesInput)

	var values []float64

	for _, valueString := range valueStrings {
		value, _ := strconv.ParseFloat(valueString, 64)
		values = append(values, value)
	}

	result := minmax(min, max, values...)

	fmt.Println(result)
}
