package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func handle_command(cmd string) {
	if cmd == "exit" {
		os.Exit(0)
	} else if after, found := strings.CutPrefix(cmd, "echo"); found {
		fmt.Print(strings.TrimSpace(after))
	} else {
		fmt.Printf("%s: command not found", cmd)
	}
	fmt.Print("\n")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s", err)
			os.Exit(1)
		}
		command = strings.TrimSpace(command)
		handle_command(command)
	}
}
