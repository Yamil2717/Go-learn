package main

import (
	"testing"
)

func equalSlices(a, b []float64) bool {
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

func TestMinMax(t *testing.T) {
	tests := []struct {
		name   string
		min    float64
		max    float64
		values []float64
		result []float64
	}{
		{"Valores fuera de los rangos Min y Max", 1, 10, []float64{2, 15, 5, -3, 10}, []float64{2, 5, 10}},
		{"Valor unico dentro del rango", 0, 100, []float64{42}, []float64{42}},
		{"Min mayor que Max", 10, 5, []float64{7, 6, 8}, []float64{}},
		{"Ningún valor dentro del rango", 1, 5, []float64{10, 20, 0}, []float64{}},
		{"Min y Max negativos", -10, -1, []float64{-5, 0, -20, -1}, []float64{-5, -1}},
		{"Min igual a Max", 5, 5, []float64{5, 4, 6, 5}, []float64{5, 5}},
		{"Incluyendo los limites del rango en los valores", 2, 8, []float64{2, 8, 3}, []float64{2, 8, 3}},
		{"Sin valores", 1, 10, []float64{}, []float64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minmax(tt.min, tt.max, tt.values...)
			if !equalSlices(got, tt.result) {
				t.Errorf("minmax(min=%v, max=%v, values=%v) [%s]: "+
					"Operacion de filtrado fallo: Se esperaba %v, se obtuvo %v",
					tt.min, tt.max, tt.values, tt.name, tt.result, got)
			}
		})
	}
}
