// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sock-Net/gows"
	"github.com/TwiN/go-color"
)

// WebsocketMessage struct
type WebSocketMessage struct {
	Type    int    `json:"type"`
	Message string `json:"data"`
	From    string `json:"from"`
}

func (c *Client) New() {
	c.Socket = gows.New(c.PrepareUrl(), true)
}

func (c *Client) Connect() error {
	err := c.Socket.Connect()
	if err != nil {
		return err
	}

	return OnReady()
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

	c.Socket.OnDisconnected = func(err error) {
		fmt.Println(color.Bold + color.Red + "Client is disconnected" + color.Reset)
		os.Exit(1)
	}
}

func OnReady() error {
	fmt.Println(color.Bold + color.Red + "Client is ready" + color.Reset)
	return nil
}
