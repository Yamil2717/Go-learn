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

Si no se pasa `-out`, se crea un archivo temporal con el patrón `md*.html` en el directorio actual (el `*` es reemplazado por un número aleatorio). El nombre del archivo generado se imprime en stdout al terminar.

### Ejemplos

```sh
# genera readme.html
go run main.go -in README.md -out readme

# genera README.md.html usando el nombre del archivo de entrada
go run main.go -in README.md
# genera md<numero>.html (archivo temporal)
go run main.go -in README.md
```
