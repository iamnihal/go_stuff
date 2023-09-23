package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	urlFlag := flag.String("u", "", "Image URL")
	fileFlag := flag.String("f", "", "File containing URLs")
	proxyFileFlag := flag.String("p", "", "File containing proxies")

	flag.Usage = func() {
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	proxyList, err := readProxyFromFile(*proxyFileFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *urlFlag == "" && *fileFlag == "" {
		fmt.Println("Usage: freepik.go -u <url> -f <file>")
		flag.Usage()
		return
	} else if *urlFlag != "" || *fileFlag != "" {
		if *urlFlag != "" {
			err := downloader(*urlFlag, []string{}, proxyList)
			if err != nil {
				fmt.Println(err)
			}
		}
		urls, err := readURLFromFile(*fileFlag)
		err = downloader("", urls, proxyList)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func readURLFromFile(f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	var urls []string
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Error opening file: %s", err)
		}

		if line != "" {
			urls = append(urls, strings.TrimSpace(line))
		}
		if err == io.EOF {
			break
		}
	}
	return urls, nil
}

func readProxyFromFile(f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %s", err)
	}

	defer file.Close()

	var proxyList []string

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Error reading file: %s", err)
		}
		if line != "" {
			proxyList = append(proxyList, strings.TrimSpace(line))
		}
		if err == io.EOF {
			break
		}
	}

	return proxyList, nil
}

func downloadEngine(u string, p string) error {
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

	fmt.Println(match[1])

	if match == nil || len(match) < 2 {
		return fmt.Errorf("Error: Unable to extract the file ID")
	}

	du := "https://www.freepik.com/download-file/" + match[1]

	proxyURLParsed, err := url.Parse(p)

	if err != nil {
		return fmt.Errorf("Error while parsing proxy URL: %s", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURLParsed),
	}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", du, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36")

	resp, err = client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error while downloding file: HTTP Status Code %v", resp.StatusCode)
	}

	defer resp.Body.Close()

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

func downloader(u string, uf []string, p []string) error {
	if len(p) == 0 {
		return fmt.Errorf("Proxy file is missing")
	}
	if u != "" && len(uf) == 0 && len(p) != 0 {
		for _, proxy := range p {
			err := downloadEngine(u, proxy)
			if err != nil {
				fmt.Println(err)
			} else {
				break
			}
		}
	} else if u == "" && len(p) != 0 {
		for _, proxy := range p {
			for _, url := range uf {
				err := downloadEngine(url, proxy)
				if err != nil {
					fmt.Println(err)
				} else {
					break
				}
			}
		}
	}
	return nil
}
