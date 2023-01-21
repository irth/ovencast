package ws

import (
	"encoding/json"
	"fmt"
)

type CommandMeta struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type Command[RequestT any, ReplyT any] struct {
	CommandMeta
	Request RequestT `json:"request"`

	wsconn *Conn
}

type RawCommand = Command[json.RawMessage, interface{}]

type FromRawable interface {
	FromRaw(r RawCommand) (Errable, error)
}

type Errable interface {
	Err(format string, args ...interface{}) error
}

func (c Command[RequestT, ReplyT]) FromRaw(r RawCommand) (Errable, error) {
	cmd := Command[RequestT, ReplyT]{
		CommandMeta: r.CommandMeta,

		wsconn: r.wsconn,
	}

	err := json.Unmarshal(r.Request, &cmd.Request)
	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	return cmd, nil
}

func (c Command[RequestT, ReplyT]) OK(reply ReplyT) error {
	response := Response[ReplyT]{
		CommandMeta: c.CommandMeta,
		OK:          true,
		Response:    reply,
	}
	return c.wsconn.SendRaw(response)
}

func (c Command[RequestT, ReplyT]) Err(format string, args ...interface{}) error {
	response := Response[ReplyT]{
		CommandMeta: c.CommandMeta,
		OK:          false,
		Error:       fmt.Sprintf(format, args...),
	}

	return c.wsconn.SendRaw(response)
}
