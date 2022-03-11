// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Sock-Net/gows"
	"github.com/TwiN/go-color"
)

// WebsocketMessage struct
type WebSocketMessage struct {
	Type    int    `json:"type,omitempty"`
	Message string `json:"data,omitempty"`
	From    string `json:"from,omitempty"`
}

// Create new socket instance
func (c *Client) New() {
	c.Socket = gows.New(c.PrepareUrl(), true)
}

// Connect to websocket
func (c *Client) Connect() {
	c.Socket.Connect()
	c.Ready()
}

// Set message handler
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

// Minimal shell for sockc
func (c *Client) Ready() {
	// Ping every 30 second
	go func() {
		time.Sleep(30 * time.Second)
		c.Socket.SendJSON(map[string]interface{}{
			"type": -1,
			"data": "",
		})
	}()

	for {
		command := HandleStdin("")

		switch strings.ToLower(command) {
		case "message":
			c.SendMessageHandler()
		}
	}
}

// Send message to other instances
func (c *Client) SendMessageHandler() {
	message := HandleStdin(color.Bold + "Message: " + color.Reset)
	websocketMessage := new(WebSocketMessage)
	websocketMessage.Type = 1
	websocketMessage.Message = message

	err := c.Socket.SendJSON(websocketMessage)
	if err != nil {
		log.Fatal(color.Bold + color.Red + "Failed to send message" + color.Reset)
	}

	fmt.Println(color.Bold + color.Green + "Message sent" + color.Reset)
}

// Handle messages from other instances
func HandleMessage(message *WebSocketMessage) {
	fmt.Println("[" + color.Bold + color.Green + "@" + message.From + color.Reset + "]: " + message.Message)
}
