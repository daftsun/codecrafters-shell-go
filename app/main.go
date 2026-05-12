package main

import (
	"bufio"
	"fmt"
	"os"
)

// var _ = fmt.Print

func main() {
	fmt.Print("$ ")
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s: command not found", command[:len(command)-1])
}
