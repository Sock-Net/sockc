// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
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
		case 2: // Send message
			HandleFile(websocketMessage)
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
		c.Socket.SendJSON(map[string]interface{}{
			"type": -1,
			"data": "",
		})
		time.Sleep(15 * time.Second)
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
		case "file":
			c.SendFileHandler(commandSplit[1:])
		case "list":
			c.ListInstancesHandler()
		}
	}
}

// Send message to other instances
func (c *Client) SendMessageHandler(args []string) {
	if len(args) == 0 {
		WriteError("Invalid argument(s)\n")
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

// Send file to other instances
func (c *Client) SendFileHandler(args []string) {
	if len(args) == 0 {
		WriteError("Invalid argument(s)\n")
		return
	}

	fileData, err := os.ReadFile(args[0])
	if err != nil {
		WriteError(err.Error() + "\n")
		return
	}

	fileName := path.Base(args[0])

	websocketMessage := new(WebSocketMessage)
	websocketMessage.Type = 2
	websocketMessage.Message = fileName + ":" + base64.StdEncoding.EncodeToString(fileData)

	err = c.Socket.SendJSON(websocketMessage)
	if err != nil {
		WriteError("Failed to send\n")
		os.Exit(1)
	}

	WriteSuccess("File sent\n")
}

// Handle files from other instances
func HandleFile(message *WebSocketMessage) {
	fileSplitted := strings.SplitN(message.Message, ":", 2)
	if len(fileSplitted) < 2 {
		return
	}

	fileName := fileSplitted[0]
	fileData, err := base64.StdEncoding.DecodeString(fileSplitted[1])
	if err != nil {
		return
	}

	os.WriteFile(fileName, fileData, os.ModeAppend)
	fmt.Println("[" + color.Bold + color.Green + "@" + message.From + color.Reset + "] sent you a file (" + color.Bold + color.Blue + fileName + color.Reset + ")")
}
