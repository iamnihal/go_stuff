package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	var url string
	files := os.Args[1:]
	for _, arg := range files {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		defer f.Close()
		input := bufio.NewScanner(f)
		for input.Scan() {
			url = input.Text()
			if strings.HasPrefix(url, "https://") != true {
				url = fmt.Sprintf("https://%v", url)
			}
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "status: %v\n", err)
				continue
			}
			defer resp.Body.Close()
			fmt.Printf("%v : [%v]\n", url, resp.Status)
		}
	}
}
