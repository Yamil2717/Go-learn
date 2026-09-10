package todo

import (
	"os"
	"testing"
)

func findTask(t *testing.T, l List, taskName string) int {
	t.Helper()
	for i, it := range l {
		if it.Task == taskName {
			return i
		}
	}
	t.Fatalf("La tarea %q no se encontro en la lista", taskName)
	return -1
}

func TestAdd(t *testing.T) {
	l := &List{}

	l.Add("tarea 1")
	if len(*l) != 1 {
		t.Fatalf("Se esperaba 1 item, la lista tiene %d", len(*l))
	}
	if i := findTask(t, *l, "tarea 1"); i != 0 {
		t.Fatalf("Se esperaba la tarea %q en el indice 0, esta en el %d", "tarea 1", i)
	}

	l.Add("tarea 2")
	if len(*l) != 2 {
		t.Fatalf("Se esperaban 2 items, la lista tiene %d", len(*l))
	}
	if i := findTask(t, *l, "tarea 2"); i != 1 {
		t.Fatalf("Se esperaba la tarea %q en el indice 1, esta en el %d", "tarea 2", i)
	}
}

func TestComplete(t *testing.T) {
	l := &List{}
	l.Add("tarea 1")
	l.Add("tarea 2")

	if err := l.Complete(1); err != nil {
		t.Fatalf("Complete fallo: %v", err)
	}

	if len(*l) != 2 {
		t.Fatalf("La longitud de la lista no deberia cambiar, tiene %d", len(*l))
	}

	i := findTask(t, *l, "tarea 2")
	if !(*l)[i].Done {
		t.Errorf("La tarea %q deberia estar completada", "tarea 2")
	}
	if (*l)[i].CompletedAt.IsZero() {
		t.Errorf("La tarea %q deberia tener CompletedAt asignado", "tarea 2")
	}

	i = findTask(t, *l, "tarea 1")
	if (*l)[i].Done {
		t.Errorf("La tarea %q no deberia estar completada", "tarea 1")
	}
	if !(*l)[i].CompletedAt.IsZero() {
		t.Errorf("La tarea %q no deberia tener CompletedAt asignado", "tarea 1")
	}

	if err := l.Complete(2); err == nil {
		t.Error("Se esperaba un error al completar con un indice invalido")
	}
}

func TestDelete(t *testing.T) {
	l := &List{}
	l.Add("tarea 1")
	l.Add("tarea 2")
	l.Add("tarea 3")

	if err := l.Delete(1); err != nil {
		t.Fatalf("Delete fallo: %v", err)
	}

	if len(*l) != 2 {
		t.Fatalf("Se esperaban 2 items, la lista tiene %d", len(*l))
	}
	for _, it := range *l {
		if it.Task == "tarea 2" {
			t.Errorf("La tarea %q deberia haber sido eliminada", "tarea 2")
		}
	}

	i := findTask(t, *l, "tarea 1")
	if i != 0 {
		t.Errorf("Se esperaba la tarea %q en el indice 0, esta en el %d", "tarea 1", i)
	}
	i = findTask(t, *l, "tarea 3")
	if i != 1 {
		t.Errorf("Se esperaba la tarea %q en el indice 1, esta en el %d", "tarea 3", i)
	}

	if err := l.Delete(5); err == nil {
		t.Error("Se esperaba un error al eliminar con un indice invalido")
	}
}

func TestSaveAndGet(t *testing.T) {
	tf, err := os.CreateTemp("", "todo")
	if err != nil {
		t.Fatalf("Error creando el archivo temporal: %v", err)
	}
	filename := tf.Name()
	if err := tf.Close(); err != nil {
		t.Fatalf("Error cerrando el archivo temporal: %v", err)
	}
	defer func() {
		if err := os.Remove(filename); err != nil {
			t.Fatalf("Error eliminando el archivo temporal %s: %v", filename, err)
		}
	}()

	original := &List{}
	original.Add("tarea 1")
	original.Add("tarea 2")
	original.Add("tarea 3")
	original.Complete(1)

	if err := original.Save(filename); err != nil {
		t.Fatalf("Save fallo: %v", err)
	}

	l2 := &List{}
	if err := l2.Get(filename); err != nil {
		t.Fatalf("Get fallo: %v", err)
	}

	if len(*l2) != len(*original) {
		t.Fatalf("Se esperaban %d items, la lista tiene %d", len(*original), len(*l2))
	}

	for i, it := range *original {
		got := (*l2)[i]
		if got.Task != it.Task {
			t.Fatalf("Se esperaba la tarea %q en el indice %d, se obtuvo %q", it.Task, i, got.Task)
		}
		if got.Done != it.Done {
			t.Errorf("La tarea %q: se esperaba Done=%v, se obtuvo Done=%v", it.Task, it.Done, got.Done)
		}
		if !got.CreatedAt.Equal(it.CreatedAt) {
			t.Errorf("La tarea %q: CreatedAt no coincide", it.Task)
		}
		if !got.CompletedAt.Equal(it.CompletedAt) {
			t.Errorf("La tarea %q: CompletedAt no coincide", it.Task)
		}
	}
}
