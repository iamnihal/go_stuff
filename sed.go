package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		search    string
		input     []string
		newStr    string
		currLine  string
		globalOpt int
		g         string
		success   bool
	)

	flag.StringVar(&search, "s", "", "s/<search>/<replace>/g")
	flag.Parse()

	params := strings.Split(search, "/")
	if params[0] != "s" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	searchStr := params[1]
	replaceStr := params[2]

	if len(params) <= 3 && params[0] == "s" || params[3] == "" {
		globalOpt = 1
	} else {
		g = params[3]
		var err error
		globalOpt, err = strconv.Atoi(g)
		if err != nil {
			fmt.Println(err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		currLine = scanner.Text()
		if strings.Contains(currLine, searchStr) && !success {
			newStr = strings.Replace(currLine, searchStr, replaceStr, globalOpt)
			success = true
			input = append(input, newStr)
		} else {
			input = append(input, currLine)
		}
	}
	for _, lines := range input {
		fmt.Println(lines)
	}
}
