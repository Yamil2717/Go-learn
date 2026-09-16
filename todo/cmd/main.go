package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"todo"
)

const todoFileName = ".todo.json"

func main() {
	complete := flag.Int("complete", -1, "Marcar una tarea como completada por su indice")
	delete := flag.Int("delete", -1, "Eliminar una tarea por su indice")

	flag.Parse()

	l := &todo.List{}

	data, err := os.ReadFile(todoFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(bytes.TrimSpace(data)) > 0 {
		if err := l.Get(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	args := flag.Args()

	switch {
	case *complete > -1:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Se completo la tarea %d\n", *complete)
	case *delete > -1:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Se elimino la tarea %d\n", *delete)
	default:
		task := strings.Join(args, " ")
		if task == "" {
			if len(*l) == 0 {
				fmt.Println("No hay tareas ingresadas")
				os.Exit(0)
			}
			for _, item := range *l {
				fmt.Println(item.Task)
			}
			os.Exit(0)
		}
		l.Add(task)
		fmt.Printf("Se agrego la tarea: %s\n", task)
	}

	if err := l.Save(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
