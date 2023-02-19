package chat

import (
	"log"

	"github.com/irth/wsrpc"
)

type MessageRequest struct {
	Content string `json:"content"`
}
type MessageCommand = wsrpc.Command[MessageRequest, any]

type Message struct {
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

func (m *Message) Type() string { return "message" }

func (c *Chat) handleMessage(state *ClientState, cmd MessageCommand) {
	if !state.ready {
		cmd.Err("you cannot send messages until you set a nickname")
		return
	}

	err := c.ch.Broadcast(state.ctx, &Message{
		Nickname: state.nick,
		Content:  cmd.Request.Content,
	})
	if err != nil {
		log.Println("broadcast error:", err)
		cmd.Err("internal server error")
		return
	}

	cmd.OK(nil)
}
