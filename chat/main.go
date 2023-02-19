package main

import (
	"context"

	"github.com/irth/ovencast/chat/internal/chat"
)

func main() {
	chat, err := chat.NewChat()
	if err != nil {
		panic(err)
	}
	go chat.Start(context.Background())

	panic(chat.Listen(":6214"))
}
