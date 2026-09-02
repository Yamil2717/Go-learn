package main

import (
	"strings"
	"testing"
)

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCapturaDeEntrada(t *testing.T) {
	tests := []struct {
		name  string
		input string
		value []string
	}{
		{"Cuando el usuario escribe exit al inicio de una linea", "exit\nhello", []string{}},
		{"Cuando el usuario escribe exit en medio de las lineas", "hello\nexit\nworld", []string{"hello"}},
		{"Cuando el usuario escribe exit al final de las lineas", "hello\nworld\nexit", []string{"hello", "world"}},
		{"Cuando el usuario escribe Exit en mayuscula y no corta", "Hello\nExit\nWorld", []string{"Hello", "Exit", "World"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := capture(strings.NewReader(tt.input))
			if !equalStringSlices(got, tt.value) {
				t.Errorf("capture(input=%q) [%s]: Operacion de captura de entrada fallo: Se esperaba %v, Se obtuvo %v",
					tt.input, tt.name, tt.value, got)
			}
		})
	}
}

func TestConteoDePalabras(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		value      int
		countLines bool
		countBytes bool
	}{
		{"Cuando el usuario escribe una sola frase con algunas palabras", "La casa es grande", 4, false, false},
		{"Cuando el usuario escribe varias frases con algunas palabras por frase", "Hola mundo. Adios gente.", 4, false, false},
		{"Cuando el usuario escribe una sola palabra", "palabra", 1, false, false},
		{"Cuando el usuario escribe una palabra compuesta como read-only", "read-only", 1, false, false},
		{"Cuando el usuario escribe multiples saltos de linea entre palabras", "uno\n\ndos", 2, false, false},
		{"Cuando el usuario escribe Exit como palabra", "Exit", 1, false, false},
		{"Cuando el usuario escribe exit como palabra", "exit", 1, false, false},
		{"Cuando el usuario escribe EXIT como palabra", "EXIT", 1, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count(strings.Split(tt.input, "\n"), tt.countLines, tt.countBytes)
			if got != tt.value {
				t.Errorf("count(input=%q, countLines=%v, countBytes=%v) [%s]: Operacion de conteo de palabras fallo: Se esperaba %d, Se obtuvo %d",
					tt.input, tt.countLines, tt.countBytes, tt.name, tt.value, got)
			}
		})
	}
}

func TestConteoDeLineas(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		value      int
		countLines bool
		countBytes bool
	}{
		{"cuando el usuario escribe una sola linea", "hola mundo", 1, true, false},
		{"cuando el usuario escribe multiples lineas sin salto de linea entre ellas", "primera linea segunda linea", 1, true, false},
		{"cuando el usuario escribe multiples lineas con saltos de linea entre ellas", "linea1\nlinea2\nlinea3", 3, true, false},
		{"cuando el usuario escribe multiples lineas con mas de un salto de linea entre ellas", "linea1\n\nlinea2", 3, true, false},
		{"cuando el usuario pasa las banderas l y b se cuentan lineas", "hola mundo", 1, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count(strings.Split(tt.input, "\n"), tt.countLines, tt.countBytes)
			if got != tt.value {
				t.Errorf("count(input=%q, countLines=%v, countBytes=%v) [%s]: Operacion de conteo de lineas fallo: Se esperaba %d, Se obtuvo %d",
					tt.input, tt.countLines, tt.countBytes, tt.name, tt.value, got)
			}
		})
	}
}

func TestConteoDeBytes(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		value      int
		countLines bool
		countBytes bool
	}{
		{"Cuando el usuario escribe una sola linea", "hola mundo", 10, false, true},
		{"Cuando el usuario escribe multiples lineas", "linea1\nlinea2", 13, false, true},
		{"Cuando el usuario escribe multiples lineas con mas de un salto de linea entre ellas", "linea1\n\nlinea2", 14, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count(strings.Split(tt.input, "\n"), tt.countLines, tt.countBytes)
			if got != tt.value {
				t.Errorf("count(input=%q, countLines=%v, countBytes=%v) [%s]: Operacion de conteo de bytes fallo: Se esperaba %d, Se obtuvo %d",
					tt.input, tt.countLines, tt.countBytes, tt.name, tt.value, got)
			}
		})
	}
}
