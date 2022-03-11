// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Sock struct
type Sock struct {
	Channel   string    `json:"channel"`
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// Get all socks (instances) from channel
func (c *Client) GetChannelSocks(channel string) []*Sock {
	resp, err := http.Get(c.PrepareHTTPUrl(""))
	if err != nil {
		log.Fatal("HTTP Error", err)
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reader Error", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var socks []*Sock
	err = json.Unmarshal(bytes, &socks)
	if err != nil {
		socks := []*Sock{}
		return socks
	}

	return socks
}
