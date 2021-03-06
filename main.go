// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/TwiN/go-color"
)

func main() {
	// Setup connection options
	host := HandleStdin(color.Bold + "(1) Host to connect: " + color.Reset)
	channel := HandleStdin(color.Bold + "(2) Channel to connect: " + color.Reset)
	id := HandleStdin(color.Bold + "(3) Id for yourself: " + color.Reset)
	token := flag.String("token", "demo", "Set connection token.")
	secureConnection := flag.Bool("secure", true, "Set secure connection (wss).")
	flag.Parse()

	client := Client{
		SecureConnection: *secureConnection,
		Host:             host,
		Channel:          channel,
		Token:            *token,
		Id:               id,
	}

	client.New()
	client.SetHandler()
	client.Connect()

	// Wait for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}

// Write text and read string until new line
func HandleStdin(text string) string {
	fmt.Print(text)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(input)
}

// Write error message to stdout
func WriteError(text string) {
	fmt.Print(color.Bold + color.Red + text + color.Reset)
}

// Write success message to stdout
func WriteSuccess(text string) {
	fmt.Print(color.Bold + color.Green + text + color.Reset)
}
