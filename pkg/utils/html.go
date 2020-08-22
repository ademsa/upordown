package utils

import (
	"golang.org/x/net/html"
	"log"
	"strings"
)

// FindNestedURLs to recursively look into html.Node for urls
func FindNestedURLs(n *html.Node, nestedUrls map[string]int, targetURL string, verbose bool) map[string]int {
	if n.Type == html.ElementNode &&
		(n.Data == "link" || n.Data == "script" || n.Data == "img" || n.Data == "a") {
		for _, v := range n.Attr {
			if v.Key == "href" || v.Key == "src" {
				url := v.Val
				if string(url[0:4]) != "http" {
					if strings.HasPrefix(url, "/") && strings.HasSuffix(targetURL, "/") {
						url = targetURL + url[1:]
					} else {
						url = targetURL + url
					}
				}
				if verbose {
					log.Println(url)
				}
				nestedUrls[url]++
			}
		}
	}

	for nc := n.FirstChild; nc != nil; nc = nc.NextSibling {
		nestedUrls = FindNestedURLs(nc, nestedUrls, targetURL, verbose)
	}

	return nestedUrls
}
