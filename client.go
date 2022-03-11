// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"github.com/Sock-Net/gows"
)

// Client struct
type Client struct {
	SecureConnection         bool
	Host, Channel, Token, Id string
	Socket                   *gows.Socket
}

// Prepare websocket url for client
func (c *Client) PrepareUrl() string {
	url := "ws"
	query := []string{}

	if c.SecureConnection {
		url += "s://"
	} else {
		url += "://"
	}

	url += c.Host + "/channel/" + c.Channel

	if c.Token != "" {
		query = append(query, "token="+c.Token)
	}

	if c.Id != "" {
		query = append(query, "id="+c.Id)
	}

	if len(query) > 0 {
		url += "?"

		for _, i := range query {
			url += i + "&"
		}

		url = url[:len(url)-1]
	}

	return url
}

// Prepare http url for client
func (c *Client) PrepareHTTPUrl(path string) string {
	url := "http"

	if c.SecureConnection {
		url += "s://"
	} else {
		url += "://"
	}

	url += c.Host + "/api/channel/" + c.Channel + path
	return url
}
