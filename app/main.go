package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = []string{"exit", "echo", "type"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, _ := reader.ReadString('\n')

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		command := fields[0]
		switch command {
		case "exit":
			return
		case "echo":
			fmt.Println(strings.Join(fields[1:], " "))
		case "type":
			type_command(fields[1:])
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func type_command(args []string) {
	if len(args) > 0 {
		for _, arg := range args {
			if slices.Contains(builtins, arg) {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else if path, _ := exec.LookPath(arg); path != "" {
				fmt.Printf("%s is %s\n", arg, path)
			} else {
				fmt.Printf("%s: not found\n", arg)
			}
		}
	}
}
