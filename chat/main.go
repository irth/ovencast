package main

import "github.com/irth/ovencast/chat/internal/chat"

func main() {
	chat, err := chat.NewChat()
	if err != nil {
		panic(err)
	}

	panic(chat.Listen(":6214"))
}
