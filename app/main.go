package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s", err)
			os.Exit(1)
		}
		command = command[:len(command)-1]
		if command == "exit" {
			os.Exit(0)
		}
		fmt.Printf("%s: command not found\n", command)
	}
}
