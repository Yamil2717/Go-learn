package main

import (
	"flag"
	"fmt"
	"os"
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
	out := flag.String("out", "", "Nombre del archivo HTML de salida, sin la extension .html")
	flag.Parse()

	if err := run(*out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(out string) error {
	if out == "" {
		return fmt.Errorf("El flag -out es obligatorio. Pasa el nombre del archivo sin la extension .html")
	}

	data := []byte(header + footer)

	if err := saveHTML(out+".html", data); err != nil {
		return err
	}

	fmt.Printf("Archivo %s.html generado\n", out)
	return nil
}

func saveHTML(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
