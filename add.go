package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var prefix string
	var suffix string
	var input string

	flag.StringVar(&prefix, "prefix", "", "Add prefix to output")
	flag.StringVar(&suffix, "suffix", "", "Add suffix to output")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = scanner.Text()
	}

	if prefix != "" && suffix == "" {
		fmt.Println(prefix, input)
	}

	if suffix != "" && prefix == "" {
		fmt.Println(input, suffix)
	}

	if suffix != "" && prefix != "" {
		fmt.Println(prefix, input, suffix)
	}

	if suffix == "" && prefix == "" {
		fmt.Println(input)
	}

}
