package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	err := build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	err = os.WriteFile(fileName, []byte{}, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file %s", fileName)
		os.Exit(1)
	}

	fmt.Println("Running tests....")
	result := m.Run()

	fmt.Println("Cleaning up....")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "New Task"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)
		fmt.Println(cmd)
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		fmt.Println(cmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf("- [ ] 0: %s", task)
		got := strings.TrimSpace(string(out))
		if !strings.Contains(got, want) {
			t.Errorf("expected %q to contain %q", got, want)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "0")
		fmt.Println(cmd)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		listCmd := exec.Command(cmdPath, "-list")
		out, err := listCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		want := fmt.Sprintf("- [X] 0: %s", task)
		got := strings.TrimSpace(string(out))
		if !strings.Contains(got, want) {
			t.Errorf("expected %q to contain %q", got, want)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		if err := exec.Command(cmdPath, "-delete", "0").Run(); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fileName, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		if err := exec.Command(cmdPath, "-task", task).Run(); err != nil {
			t.Fatal(err)
		}

		if err := exec.Command(cmdPath, "-delete", "0").Run(); err != nil {
			t.Fatal(err)
		}

		listCmd := exec.Command(cmdPath, "-list")
		out, err := listCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		if out := strings.TrimSpace(string(out)); out != "" {
			t.Errorf("expected no tasks, got %q instead", out)
		}
	})
}
