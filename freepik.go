package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	urlFlag := flag.String("u", "", "Image URL")

	flag.Parse()

	if *urlFlag == "" {
		println("URL is missing")
	} else {
		err := downloadImage(*urlFlag)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func downloadImage(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: Unexpected status code %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response: %s", err)
	}
	pattern := regexp.MustCompile(`download-file/([0-9]+)`)
	match := pattern.FindStringSubmatch(string(body))

	if len(match) < 2 {
		return fmt.Errorf("Error: Unable to extract the file ID")
	}
	resp, err = http.Get("https://www.freepik.com/download-file/" + match[1])
	if err != nil {
		return fmt.Errorf("Error while downloding file: %s", err)
	}

	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	file, err := os.Create(match[1] + ".png")
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Error saving file: %s", err)
	}
	return nil
}
