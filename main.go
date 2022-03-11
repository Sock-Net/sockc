// Copyright (c) 2022 aiocat
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
}

func HandleMessage(message *WebSocketMessage) {
	fmt.Println(message.Message)
}

func HandleStdin(text string) string {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
}
