// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"encoding/json"
	"fmt"
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
	for _, sock := range c.GetChannelSocks(c.Channel) {
		if sock.Id == c.Id {
			WriteError("ID already taken\n")
			os.Exit(1)
		}
	}

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
		WriteError("Client disconnected\n")
		os.Exit(1)
	}

	c.Socket.OnConnectError = func(err error) {
		WriteError("Client is unable to connect\n")
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
	WriteSuccess("Ready!\n")

	for {
		command := HandleStdin("")
		commandSplit := strings.SplitN(command, " ", 2)

		if commandSplit == nil {
			WriteError("Invalid command format\n")
			continue
		}

		switch strings.ToLower(commandSplit[0]) {
		case "message":
			c.SendMessageHandler(commandSplit[1:])
		case "list":
			c.ListInstancesHandler()
		}
	}
}

// Send message to other instances
func (c *Client) SendMessageHandler(args []string) {
	if len(args) == 0 {
		WriteError("Invalid argument\n")
		return
	}

	websocketMessage := new(WebSocketMessage)
	websocketMessage.Type = 1
	websocketMessage.Message = args[0]

	err := c.Socket.SendJSON(websocketMessage)
	if err != nil {
		WriteError("Failed to send\n")
		os.Exit(1)
	}

	WriteSuccess("Message sent\n")
}

// Handle messages from other instances
func HandleMessage(message *WebSocketMessage) {
	fmt.Println("[" + color.Bold + color.Green + "@" + message.From + color.Reset + "]: " + message.Message)
}

// List all instances in channel
func (c *Client) ListInstancesHandler() {
	instances := c.GetChannelSocks(c.Channel)
	if len(instances) == 0 {
		WriteError("No instance found\n")
		return
	}

	for _, instance := range instances {
		fmt.Println("[" + color.Bold + color.Green + instance.Id + color.Reset + ":" + color.Bold + color.Yellow + instance.Channel + color.Reset + "] at " + color.Blue + instance.CreatedAt.String() + color.Reset)
	}
}
