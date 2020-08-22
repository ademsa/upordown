package utils

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

// Worker for reading and executing jobs
func Worker(id int, stopWorkers <-chan bool, jobs <-chan map[string]interface{}, results chan<- map[string]interface{}) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	defer httpClient.CloseIdleConnections()

	for {
		select {
		case <-stopWorkers:
			return
		case j := <-jobs:
			r := j

			if j == nil {
				return
			}

			st := time.Now()
			response, err := httpClient.Get(r["url"].(string))
			et := time.Now()
			if err != nil {
				r["err"] = err
				r["rt"] = time.Duration(0).Milliseconds()
				if strings.Contains(err.Error(), "Client.Timeout") {
					r["statusCode"] = 0
				} else {
					r["statusCode"] = 404
				}
			} else {
				r["rt"] = et.Sub(st).Milliseconds()
				r["statusCode"] = response.StatusCode
			}

			results <- r
		case <-time.After(15 * time.Second):
			break
		}
	}
}

// AddJobs to add jobs to jobs channel and to send preliminary nested url data to websocket client
func AddJobs(jobs chan<- map[string]interface{}, nestedUrls map[string]int, c *websocket.Conn) error {
	idx := 1
	for url, urlCount := range nestedUrls {
		j := map[string]interface{}{
			"id":       idx,
			"url":      url,
			"urlCount": urlCount,
		}
		jobs <- j
		idx++
		if c != nil {
			err := SendNestedURLMessage(c, "nested-url-in-progress", j["id"].(int), j["url"].(string), j["urlCount"].(int), 0, 0)
			if err != nil {
				close(jobs)
				return err
			}
		}
	}
	close(jobs)
	return nil
}

// ProcessResults to process results from results channel and to send nested url results data to websocket client
func ProcessResults(results <-chan map[string]interface{}, stopResults <-chan bool, nestedUrlsCount int, c *websocket.Conn, verbose bool) error {
	for i := 0; i < nestedUrlsCount; i++ {
		select {
		case <-stopResults:
			return nil
		case r := <-results:
			if verbose && r["rt"].(int64) == 0 {
				log.Println(r["id"], r["statusCode"], r["url"])
			} else if verbose && r["rt"].(int64) > 0 {
				log.Println(r["id"], r["statusCode"], r["url"], r["url"], ">", r["rt"], "ms")
			}

			if c != nil {
				err := SendNestedURLMessage(c, "nested-url-completed", r["id"].(int), r["url"].(string), r["urlCount"].(int), r["rt"].(int64), r["statusCode"].(int))
				if err != nil {
					return err
				}
			}
		case <-time.After(15 * time.Second):
			return nil
		}
	}
	return nil
}
