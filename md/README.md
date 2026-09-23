# Markdown Preview Tool

Programa de consola que toma un archivo markdown y genera un archivo HTML con el resultado.

## Uso

```sh
go run main.go -in <archivo.md> [-out <nombre>]
```

### Flags

|  Flag  | Obligatorio | Descripción |
|--------|-------------|-------------|
| `-in`  | Sí | Ruta al archivo markdown de entrada |
| `-out` | No | Nombre del archivo HTML de salida, sin la extensión `.html` |

Si no se pasa `-out`, el archivo de salida se nombra a partir del archivo de entrada: `-in README.md` genera `README.md.html`.

### Ejemplos

```sh
# genera readme.html
go run main.go -in README.md -out readme

# genera README.md.html usando el nombre del archivo de entrada
go run main.go -in README.md
```
