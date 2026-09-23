package main

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Fatal(err)
		}
	}()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	t.Run("Success", func(t *testing.T) {
		if err := run("index"); err != nil {
			t.Fatalf("run fallo: %v", err)
		}

		got, err := os.ReadFile("index.html")
		if err != nil {
			t.Fatalf("No se pudo leer index.html: %v", err)
		}

		expected := []byte(header + footer)
		if string(got) != string(expected) {
			t.Errorf("Se esperaba %q, se obtuvo %q", expected, got)
		}
	})

	t.Run("MissingOutFlag", func(t *testing.T) {
		if err := run(""); err == nil {
			t.Error("Se esperaba un error al no pasar el flag -out")
		}
	})
}
