package ws

type messageWrapper struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type Message interface {
	Type() string
}

func (c *Conn) SendMessage(m Message) error {
	return c.SendRaw(messageWrapper{
		Type:    m.Type(),
		Message: m,
	})
}
