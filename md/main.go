package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

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
	out := flag.String("out", "", "Nombre del archivo HTML de salida; en blanco se usa un archivo temporal")
	flag.Parse()

	if err := run(os.Stdout, *in, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(writer io.Writer, in string, out string) error {
	if in == "" {
		return fmt.Errorf("El flag -in es obligatorio. Pasa la ruta al archivo markdown")
	}

	data, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("No se pudo leer el archivo %s: %w", in, err)
	}

	var outname string
	if out == "" {
		file, err := os.CreateTemp(".", "md*.html")
		if err != nil {
			return fmt.Errorf("No se pudo crear el archivo temporal: %w", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("No se pudo cerrar el archivo temporal: %w", err)
		}
		outname = file.Name()
	} else {
		outname = out + ".html"
	}

	body := parseContent(data)

	if err := saveHTML(outname, body); err != nil {
		return err
	}

	fmt.Fprintln(writer, outname)
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
