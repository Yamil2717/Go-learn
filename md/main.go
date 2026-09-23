package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const header = `<!DOCTYPE html>
  <html>
    <head>
      <meta http-equiv="content-type" content="text/html; charset=utf-8" />
      <title>Markdown Preview Tool</title>
    </head>
    <body>
`

const footer = `
    </body>
  </html>
`

func main() {
	in := flag.String("in", "", "Ruta al archivo markdown de entrada")
	out := flag.String("out", "", "Nombre del archivo HTML de salida, sin la extension .html")
	flag.Parse()

	if err := run(*in, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(in string, out string) error {
	if in == "" {
		return fmt.Errorf("El flag -in es obligatorio. Pasa la ruta al archivo markdown")
	}

	data, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("No se pudo leer el archivo %s: %w", in, err)
	}

	if out == "" {
		out = filepath.Base(in)
	}

	body := parseContent(data)

	if err := saveHTML(out+".html", body); err != nil {
		return err
	}

	fmt.Printf("Archivo %s.html generado\n", out)
	return nil
}

func parseContent(input []byte) []byte {
	normalized := bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n"))
	output := blackfriday.Run(normalized)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)
	return append([]byte(header), append(body, footer...)...)
}

func saveHTML(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
