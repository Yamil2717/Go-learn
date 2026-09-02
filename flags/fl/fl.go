package fl

import "os"

type Flag struct {
	value       bool
	description string
	cmd         string
}

var flags = make(map[string]*Flag)

func Parse() {
	args := os.Args[1:]

	for _, arg := range args {
		if flag, exists := flags[arg]; exists {
			flag.value = true
		}
	}
}

func Bool(cmd string, value bool, description string) *bool {
	flag := &Flag{
		cmd:         cmd,
		value:       value,
		description: description,
	}

	flags[cmd] = flag

	return &flag.value
}
