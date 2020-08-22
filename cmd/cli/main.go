package main

import (
	"flag"
	"log"
	"upordown/pkg/utils"
)

func main() {
	target := flag.String("target", "", "Target URL")
	now := flag.Int("workers", 10, "Number of workers")

	flag.Parse()

	if *target == "" {
		log.Println("[Input] Validation Error", "Target URL cannot be empty.")
		return
	}
	if *now < 1 {
		log.Println("[Input] Validation Error", "At least one worker needs to be enabled.")
		return
	}

	log.Println("UpOrDown")
	log.Println("[Settings] Number of workers:", *now)
	log.Println("Target URL:", *target)

	baseURL, statusCode, rt, htmlTree, err := utils.CheckTargetURL(*target)
	if err != nil {
		log.Println("Down! > Error:", err)
	} else if htmlTree == nil {
		log.Println("Status Code:", statusCode)
		log.Println("Response Time:", rt, "ms")
		log.Println("No nested urls.")
		return
	} else {
		log.Println("Status Code:", statusCode)
		log.Println("Response Time:", rt, "ms")
	}

	log.Println("Looking for nested URLs")

	nestedUrls := make(map[string]int)
	utils.FindNestedURLs(htmlTree, nestedUrls, baseURL, true)

	if len(nestedUrls) == 0 {
		log.Println("No nested urls.")
		return
	}

	log.Println("Checking nested URLs")

	jobs := make(chan map[string]interface{}, len(nestedUrls))
	results := make(chan map[string]interface{}, len(nestedUrls))
	stopWorkers := make(chan bool)
	stopResults := make(chan bool)

	for i := 0; i < *now; i++ {
		go utils.Worker(i, stopWorkers, jobs, results)
	}

	err = utils.AddJobs(jobs, nestedUrls, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.ProcessResults(results, stopResults, len(nestedUrls), nil, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Completed checks.")
}
