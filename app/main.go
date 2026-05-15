package main

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
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

func readCommand(reader readline.Instance) ([]string, error) {
	line, err := reader.Readline()
	if err != nil {
		return nil, err
	}

	fields, _ := shlex.Split(strings.TrimSpace(string(line)))
	if len(fields) == 0 {
		return nil, errors.New("no command entered")
	}
	return fields, nil
}

func handleRedirection(fields []string) (io.Writer, io.Writer, []string) {
	stdout := os.Stdout
	stderr := os.Stderr
	n := len(fields)

	if n > 2 {
		check_field := fields[n-2]
		if strings.Contains(check_field, ">") {
			switch check_field {
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
	return stdout, stderr, fields
}

func evalCommand(fields []string) {
	stdout, stderr, fields := handleRedirection(fields)

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

type RingingAutocompleter struct {
	handler readline.AutoCompleter
}

func (m *RingingAutocompleter) Do(line []rune, pos int) ([][]rune, int) {
	newLine, length := m.handler.Do(line, pos)

	if length == 0 && len(line) > 0 {
		fmt.Print("\a")
	}

	return newLine, length
}

func main() {

	builtins = map[string]func(cmdInputs){
		"exit": func(cmdInputs) { os.Exit(0) },
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}

	var completer []readline.PrefixCompleterInterface
	for k := range maps.Keys(builtins) {
		completer = append(completer, readline.PcItem(k))
	}
	paths := filepath.SplitList(os.Getenv("PATH"))
	for _, path := range paths {
		files, _ := os.ReadDir(path)
		for _, file := range files {
			finfo, _ := file.Info()
			if !finfo.IsDir() && finfo.Mode().Perm()&0111 != 0 {
				completer = append(completer, readline.PcItem(finfo.Name()))
			}
		}
	}

	reader, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    &RingingAutocompleter{readline.NewPrefixCompleter(completer...)},
		InterruptPrompt: "^C",
		HistoryFile:     "/tmp/shell_hist.txt",
		HistoryLimit:    10,
	})
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		return
	}
	defer reader.Close()

	for {
		fields, err := readCommand(*reader)
		if err != nil {
			continue
		}
		evalCommand(fields)
	}
}
