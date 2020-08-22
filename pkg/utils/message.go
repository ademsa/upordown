package utils

import (
	"github.com/gorilla/websocket"
)

// ProcessMessage to process websocket message
func ProcessMessage(now int, stopWorkers <-chan bool, stopResults <-chan bool, c *websocket.Conn, m map[string]interface{}) error {
	targetURL := string(m["url"].(string))

	baseURL, statusCode, rt, htmlTree, err := CheckTargetURL(targetURL)
	statusMsg := ""
	if err != nil {
		statusMsg = err.Error()
	} else if htmlTree == nil {
		statusMsg = "HTML tree empty"
	}

	err = SendTargetURLMessage(c, statusCode, statusMsg, rt)
	if err != nil {
		return err
	}

	if htmlTree == nil {
		err = SendCompletedMessage(c)
		if err != nil {
			return err
		}
		return nil
	}

	nestedUrls := make(map[string]int)
	nestedUrls = FindNestedURLs(htmlTree, nestedUrls, baseURL, false)

	jobs := make(chan map[string]interface{}, len(nestedUrls))
	results := make(chan map[string]interface{}, len(nestedUrls))

	for i := 0; i < now; i++ {
		go Worker(i, stopWorkers, jobs, results)
	}

	err = AddJobs(jobs, nestedUrls, c)
	if err != nil {
		return err
	}

	err = ProcessResults(results, stopResults, len(nestedUrls), c, false)
	if err != nil {
		return err
	}

	err = SendCompletedMessage(c)
	if err != nil {
		return err
	}

	return nil
}
