package main

import (
	"bytes"
	"os"
	"path/filepath"
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
		inPath := filepath.Join(oldWd, "testdata", "README.md")
		if err := run(inPath, "index"); err != nil {
			t.Fatalf("run fallo: %v", err)
		}

		got, err := os.ReadFile("index.html")
		if err != nil {
			t.Fatalf("No se pudo leer index.html: %v", err)
		}

		golden, err := os.ReadFile(filepath.Join(oldWd, "testdata", "README.md.golden"))
		if err != nil {
			t.Fatalf("No se pudo leer el golden file: %v", err)
		}

		if !bytes.Equal(got, golden) {
			t.Errorf("Se esperaba %q, se obtuvo %q", golden, got)
		}
	})

	t.Run("MissingInFlag", func(t *testing.T) {
		if err := run("", "index"); err == nil {
			t.Error("Se esperaba un error al no pasar el flag -in")
		}
	})
}

func TestParseContent(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "README.md"))
	if err != nil {
		t.Fatalf("No se pudo leer el archivo markdown de prueba: %v", err)
	}

	golden, err := os.ReadFile(filepath.Join("testdata", "README.md.golden"))
	if err != nil {
		t.Fatalf("No se pudo leer el golden file: %v", err)
	}

	got := parseContent(data)

	if !bytes.Equal(got, golden) {
		t.Errorf("Se esperaba %q, se obtuvo %q", golden, got)
	}
}
