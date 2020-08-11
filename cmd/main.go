package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"time"
)

func main() {
	url := flag.String("url", "", "Target URL")
	now := flag.Int("workers", 3, "Number of workers")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")

	flag.Parse()

	if *verbose {
		fmt.Println("UpOrDown")
		fmt.Println("Target URL:", *url)
	}

	htmlTree := checkTargetURL(*url, *verbose)
	if htmlTree == nil {
		return
	}

	nestedUrls := new([]map[string]interface{})
	findNestedURLs(htmlTree, nestedUrls)

	jobs := make(chan map[string]interface{}, len(*nestedUrls))
	results := make(chan map[string]interface{}, len(*nestedUrls))

	for i := 0; i < *now; i++ {
		go worker(jobs, results, *verbose)
	}

	fmt.Println("Checking nested urls next.")
	for _, nu := range *nestedUrls {
		if string(nu["url"].(string)[0]) == "/" {
			nu["url"] = *url + nu["url"].(string)
		}
		jobs <- nu
	}
	close(jobs)

	for i := 0; i < len(*nestedUrls); i++ {
		select {
		case r := <-results:
			if *verbose {
				fmt.Println(r["statusCode"], r["url"])
			}
		case <-time.After(2000 * time.Millisecond):
			return
		}
	}

	if *verbose {
		fmt.Println("Completed checks.")
	}
}

func checkTargetURL(url string, verbose bool) *html.Node {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	defer httpClient.CloseIdleConnections()

	response, err := httpClient.Get(url)
	if err != nil {
		if verbose {
			fmt.Println("Down! > Error:", err)
		}
		return nil
	}

	defer response.Body.Close()

	if verbose {
		if response.StatusCode != 200 {
			fmt.Println("Probably up > ", response.StatusCode)
			return nil
		}
		fmt.Println("It loads fine.")
	}

	htmlTree, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return htmlTree
}

func findNestedURLs(n *html.Node, nestedUrls *[]map[string]interface{}) {
	if n.Type == html.ElementNode &&
		(n.Data == "link" || n.Data == "script" || n.Data == "img" || n.Data == "a") {
		for _, v := range n.Attr {
			if v.Key == "href" || v.Key == "src" {
				*nestedUrls = append(*nestedUrls, map[string]interface{}{
					"url":  v.Val,
					"type": v.Key,
				})
			}
		}
	}

	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		findNestedURLs(nn, nestedUrls)
	}
}

func worker(jobs <-chan map[string]interface{}, results chan<- map[string]interface{}, verbose bool) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	defer httpClient.CloseIdleConnections()

	for job := range jobs {
		result := job

		if verbose {
			fmt.Println(">>>", result["url"])
		}

		response, err := httpClient.Get(result["url"].(string))

		if err != nil {
			result["err"] = err
			result["statusCode"] = 404
		} else {
			result["statusCode"] = response.StatusCode
		}

		results <- result
	}
}
