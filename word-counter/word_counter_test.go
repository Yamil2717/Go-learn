package main

import (
	"strings"
	"testing"
)

func TestConteoDePalabras(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		value      int
		countLines bool
	}{
		{"Cuando el usuario escribe una sola frase con algunas palabras", "La casa es grande", 4, false},
		{"Cuando el usuario escribe varias frases con algunas palabras por frase", "Hola mundo. Adios gente.", 4, false},
		{"Cuando el usuario escribe una sola palabra", "palabra", 1, false},
		{"Cuando el usuario escribe una palabra compuesta como read-only", "read-only", 1, false},
		{"Cuando el usuario escribe multiples saltos de linea entre palabras", "uno\n\ndos", 2, false},
		{"Cuando el usuario escribe Exit como palabra", "Exit", 1, false},
		{"Cuando el usuario escribe exit y corta el conteo", "exit", 0, false},
		{"Cuando el usuario escribe EXIT como palabra", "EXIT", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count(strings.NewReader(tt.input), tt.countLines)
			if got != tt.value {
				t.Errorf("count(input=%q, countLines=%v) [%s]: Operacion de conteo de palabras fallo: Se esperaba %d, Se obtuvo %d",
					tt.input, tt.countLines, tt.name, tt.value, got)
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
	}{
		{"cuando el usuario escribe una sola linea", "hola mundo", 1, true},
		{"cuando el usuario escribe multiples lineas sin salto de linea entre ellas", "primera linea segunda linea", 1, true},
		{"cuando el usuario escribe multiples lineas con saltos de linea entre ellas", "linea1\nlinea2\nlinea3", 3, true},
		{"cuando el usuario escribe exit al inicio de una linea", "exit\nhello", 0, true},
		{"cuando el usuario escribe exit en medio de las lineas", "hello\nexit\nworld", 1, true},
		{"cuando el usuario escribe exit al final de las lineas", "hello\nworld\nexit", 2, true},
		{"cuando el usuario escribe Exit en mayuscula y no corta", "Hello\nExit\nWorld", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count(strings.NewReader(tt.input), tt.countLines)
			if got != tt.value {
				t.Errorf("count(input=%q, countLines=%v) [%s]: Operacion de conteo de lineas fallo: Se esperaba %d, Se obtuvo %d",
					tt.input, tt.countLines, tt.name, tt.value, got)
			}
		})
	}
}
