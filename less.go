package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		lines      []string
		numOfLines int
	)

	flag.IntVar(&numOfLines, "n", 10, "Number of lines to read")
	flag.Parse()

	if numOfLines < 0 {
		fmt.Println("Number of lines cannot be negative!")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := len(lines) - 1; i >= len(lines)-numOfLines; i-- {
		fmt.Println(lines[i])
	}
}
