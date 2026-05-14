package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

type cmdInputs struct {
	writer io.Writer
	err    io.Writer
	args   []string
}

var builtins map[string]func(cmdInputs)

func handleEcho(args cmdInputs) {
	fmt.Fprintln(args.writer, strings.Join(args.args, " "))
}

func handleType(args cmdInputs) {
	if _, ok := builtins[args.args[0]]; ok {
		fmt.Printf("%s is a shell builtin\n", args.args[0])
	} else if path, err := exec.LookPath(args.args[0]); err == nil {
		fmt.Fprintf(args.writer, "%s is %s\n", args.args[0], path)
	} else {
		fmt.Printf("%s: not found\n", args.args[0])
	}
}

func handlePwd(args cmdInputs) {
	pwd, _ := os.Getwd()
	fmt.Fprintln(args.writer, pwd)
}

func handleCd(args cmdInputs) {
	if args.args[0] == "~" {
		home_dir, _ := os.UserHomeDir()
		os.Chdir(home_dir)
	} else if _, err := os.Stat(args.args[0]); err != nil {
		fmt.Fprintf(args.err, "cd: %s: No such file or directory\n", args.args[0])
	} else {
		os.Chdir(args.args[0])
	}
}

func readCommand(reader bufio.Reader) ([]string, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	fields, _ := shlex.Split(strings.TrimSpace(string(line)))
	if len(fields) == 0 {
		return nil, errors.New("no command entered")
	}
	return fields, nil
}

func evalCommand(fields []string) {
	n := len(fields)
	stdout := os.Stdout
	stderr := os.Stderr

	if n > 2 {
		if strings.Contains(fields[n-2], ">") {
			switch fields[n-2] {
			case ">", "1>":
				stdout, _ = os.Create(fields[n-1])
			case ">>", "1>>":
				stdout, _ = os.OpenFile(fields[n-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			case "2>":
				stderr, _ = os.Create(fields[n-1])
			case "2>>":
				stderr, _ = os.OpenFile(fields[n-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			}
			fields = fields[:n-2]
		}
	}

	command := fields[0]
	if cmdFunc, ok := builtins[command]; ok {
		cmdFunc(cmdInputs{stdout, stderr, fields[1:]})
	} else if _, err := exec.LookPath(command); err == nil {
		cmd := exec.Command(command, fields[1:]...)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		cmd.Run()
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

func main() {

	builtins = map[string]func(cmdInputs){
		"exit": func(cmdInputs) { os.Exit(0) },
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		fields, err := readCommand(*reader)
		if err != nil {
			continue
		}
		evalCommand(fields)
	}
}
