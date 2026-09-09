package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	item := item{
		Task:      task,
		Done:      false,
		CreatedAt: time.Now(),
	}

	*l = append(*l, item)
}

func (l *List) Complete(index int) error {
	if index < 0 || index >= len(*l) {
		return fmt.Errorf("Index invalido %d. La lista tiene %d items", index, len(*l))
	}

	(*l)[index].Done = true
	(*l)[index].CompletedAt = time.Now()

	return nil
}

func (l *List) Delete(index int) error {
	if index < 0 || index >= len(*l) {
		return fmt.Errorf("Index invalido %d. La lista tiene %d items", index, len(*l))
	}

	*l = append((*l)[:index], (*l)[index+1:]...)

	return nil
}

func (l *List) Save(filename string) error {
	jsonData, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("No se pudo convertir la lista a json: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("No se pudo guardar el archivo %s: %w", filename, err)
	}

	return nil
}

func (l *List) Get(filename string) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("El archivo %s no existe: %w", filename, err)
		}
		return fmt.Errorf("No se pudo leer el archivo %s: %w", filename, err)
	}

	if len(jsonData) == 0 {
		return errors.New("El archivo esta vacio, no hay nada que decodificar")
	}

	if err := json.Unmarshal(jsonData, l); err != nil {
		return fmt.Errorf("No se pudo decodificar el json del archivo %s: %w", filename, err)
	}

	return nil
}
