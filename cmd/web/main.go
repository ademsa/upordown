package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"upordown/pkg/utils"
)

var upgrader = utils.GetUpgrader()
var now *int

func main() {
	address := flag.String("address", "localhost:8080", "Websocket Address")
	now = flag.Int("workers", 10, "Number of workers")

	flag.Parse()

	if *address == "" {
		log.Println("[Input] Validation Error", "Websocket Address cannot be empty.")
		return
	}
	if *now < 1 {
		log.Println("[Input] Validation Error", "At least one worker needs to be enabled.")
		return
	}

	http.HandleFunc("/", ws)

	log.Fatal(http.ListenAndServe(*address, nil))
}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	messages := make(chan map[string]interface{})
	stopWorkers := make(chan bool, *now)
	stopResults := make(chan bool, 1)

	go func() {
		for m := range messages {
			utils.ProcessMessage(*now, stopWorkers, stopResults, c, m)
		}
	}()

	for {
		mt, m, err := c.ReadMessage()
		if err != nil {
			wsCleanUp(stopWorkers, stopResults, messages)
			break
		}

		if mt == websocket.TextMessage {
			var data map[string]interface{}
			json.Unmarshal(m, &data)
			typeVal, typeOK := data["type"].(string)
			if typeOK && string(typeVal) == "target-url-request" {
				messages <- data
			} else if typeOK && string(typeVal) == "target-url-cancel-request" {
				wsCleanUp(stopWorkers, stopResults, nil)
			}
		} else if mt == websocket.CloseMessage {
			wsCleanUp(stopWorkers, stopResults, messages)
			break
		}
	}
}

func wsCleanUp(stopWorkers chan<- bool, stopResults chan<- bool, messages chan map[string]interface{}) {
	for i := 0; i < *now; i++ {
		stopWorkers <- true
	}

	stopResults <- true

	if messages != nil {
		close(messages)
	}
}
