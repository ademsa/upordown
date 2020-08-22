package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

// GetUpgrader to get websocket Upgrader
func GetUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// SendNestedURLMessage to send nested url message to websocket client
func SendNestedURLMessage(c *websocket.Conn, dataType string, id int, url string, urlCount int, rt int64, statusCode int) error {
	data, err := json.Marshal(map[string]interface{}{
		"type":       dataType,
		"id":         id,
		"url":        url,
		"urlCount":   urlCount,
		"rt":         rt,
		"statusCode": statusCode,
	})
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}

// SendTargetURLMessage to send target url message to websocket client
func SendTargetURLMessage(c *websocket.Conn, statusCode int, statusMsg string, rt int64) error {
	data, err := json.Marshal(map[string]interface{}{
		"type":       "target-url-result",
		"statusCode": statusCode,
		"statusMsg":  statusMsg,
		"rt":         rt,
	})
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}

// SendCompletedMessage to send completed message to websocket client
func SendCompletedMessage(c *websocket.Conn) error {
	data, err := json.Marshal(map[string]interface{}{
		"type": "completed",
	})
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}
