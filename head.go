package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		input      []string
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
		input = append(input, scanner.Text())
	}

	if numOfLines == 0 || numOfLines > len(input) {
		numOfLines = len(input)
	}

	for i := 0; i < numOfLines; i++ {
		fmt.Println(input[i])
	}
}
