package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = []string{"exit", "echo", "type", "pwd"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		command := fields[0]
		if command == "exit" {
			return
		} else if command == "echo" {
			fmt.Println(strings.Join(fields[1:], " "))
		} else if command == "type" {
			handleType(fields[1])
		} else if command == "pwd" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Print(err)
			}
			fmt.Printf("%s\n", cwd)
		} else if _, err := exec.LookPath(command); err == nil {
			cmd := exec.Command(command, fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func handleType(command string) {
	if slices.Contains(builtins, command) {
		fmt.Printf("%s is a shell builtin\n", command)
	} else if path, err := exec.LookPath(command); err == nil {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}
