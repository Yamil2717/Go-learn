package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func normalizeLF(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

func TestRunUsingOutFlag(t *testing.T) {
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

	var buf bytes.Buffer
	inPath := filepath.Join(oldWd, "testdata", "README.md")
	if err := run(&buf, inPath, "index"); err != nil {
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

	if !bytes.Equal(normalizeLF(got), normalizeLF(golden)) {
		t.Errorf("Se esperaba %q, se obtuvo %q", golden, got)
	}
}

func TestRunWithoutOutFlag(t *testing.T) {
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

	var buf bytes.Buffer
	inPath := filepath.Join(oldWd, "testdata", "README.md")
	if err := run(&buf, inPath, ""); err != nil {
		t.Fatalf("run fallo: %v", err)
	}

	outfile := strings.TrimSpace(buf.String())

	if !strings.HasPrefix(filepath.Base(outfile), "md") || !strings.HasSuffix(outfile, ".html") {
		t.Fatalf("Se esperaba un nombre md*.html, se obtuvo %q", outfile)
	}

	got, err := os.ReadFile(outfile)
	if err != nil {
		t.Fatalf("No se pudo leer %s: %v", outfile, err)
	}

	golden, err := os.ReadFile(filepath.Join(oldWd, "testdata", "README.md.golden"))
	if err != nil {
		t.Fatalf("No se pudo leer el golden file: %v", err)
	}

	if !bytes.Equal(normalizeLF(got), normalizeLF(golden)) {
		t.Errorf("Se esperaba %q, se obtuvo %q", golden, got)
	}
}

func TestRunMissingInFlag(t *testing.T) {
	var buf bytes.Buffer
	if err := run(&buf, "", "index"); err == nil {
		t.Error("Se esperaba un error al no pasar el flag -in")
	}
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

	if !bytes.Equal(normalizeLF(got), normalizeLF(golden)) {
		t.Errorf("Se esperaba %q, se obtuvo %q", golden, got)
	}
}
