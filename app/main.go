package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = []string{"exit", "echo", "type", "pwd", "cd"}

func main() {

	cmdFuncMap := map[string]func(shellArgs []string){
		"exit": func(shellArgs []string) { os.Exit(0) },
		"echo": func(shellArgs []string) { fmt.Println(strings.Join(shellArgs, " ")) },
		"type": handleType,
		"pwd": func(shellArgs []string) {
			pwd, err := os.Getwd()
			if err == nil {
				fmt.Println(pwd)
			}
		},
		"cd": handleCd,
	}

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
		if cmdFunc, ok := cmdFuncMap[command]; ok {
			cmdFunc(fields[1:])
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

func handleType(shellArgs []string) {
	for _, command := range shellArgs {
		if slices.Contains(builtins, command) {
			fmt.Printf("%s is a shell builtin\n", command)
		} else if path, err := exec.LookPath(command); err == nil {
			fmt.Printf("%s is %s\n", command, path)
		} else {
			fmt.Printf("%s: not found\n", command)
		}
	}
}

func handleCd(shellArgs []string) {
	if _, err := os.Stat(shellArgs[0]); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", shellArgs[0])
	} else {
		os.Chdir(shellArgs[0])
	}
}
