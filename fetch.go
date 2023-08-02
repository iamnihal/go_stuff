package main

import (
	"bufio"
	"fmt"
	"io"
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
				fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
				continue
			}
			b, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
				continue
			}
			fmt.Printf("%s : [%s]\n", url, b)
		}
	}
}
