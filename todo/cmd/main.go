package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"todo"
)

const todoFileName = ".todo.json"

func main() {
	list := flag.Bool("list", false, "Listar las tareas pendientes")
	task := flag.String("task", "", "Agregar una tarea nueva")
	complete := flag.Int("complete", -1, "Marcar como completada la tarea con el indice dado")
	delete := flag.Int("delete", -1, "Eliminar la tarea con el indice dado")

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

	switch {
	case *list:
		for _, task := range *l {
			if task.Done {
				continue
			}
			fmt.Printf("Title: %s, Done: %t, CreatedAt: %s, CompletedAt: %s\n", task.Task, task.Done, task.CreatedAt, task.CompletedAt)
		}
	case *complete > -1:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > -1:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "No se especifico ninguna operacion. Usa -help para ver los comandos disponibles")
		os.Exit(1)
	}
}
