// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"encoding/json"

	"github.com/Sock-Net/gows"
)

// WebsocketMessage struct
type WebSocketMessage struct {
	Type    int    `json:"type"`
	Message string `json:"data"`
	From    string `json:"from"`
}

func (c *Client) Connect() error {
	c.Socket = gows.New(c.PrepareUrl(), true)

	err := c.Socket.Connect()
	if err != nil {
		return err
	}

	return c.OnReady()
}

func (c *Client) SetHandler() {
	c.Socket.OnTextMessage = func(message string) {
		websocketMessage := new(WebSocketMessage)
		err := json.Unmarshal([]byte(message), websocketMessage)
		if err != nil {
			return
		}

		switch websocketMessage.Type {
		case 1: // Send message
			HandleMessage(websocketMessage)
		}
	}
}
