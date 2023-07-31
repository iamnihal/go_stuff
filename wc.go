package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = scanner.Text()
	}
	fmt.Println("words:", len(strings.Fields(input)))
	fmt.Println("characters:", len(input))
}
