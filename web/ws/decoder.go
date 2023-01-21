package ws

import "fmt"

type CommandPalette map[string]FromRawable

func (c *CommandDecoder) Decode() (Errable, error) {
	rawCommand := RawCommand{}

	err := c.wsconn.RecvRaw(&rawCommand)
	if err != nil {
		return nil, err
	}

	rawCommand.wsconn = c.wsconn

	cmdType, ok := c.palette[rawCommand.Command]
	if !ok {
		rawCommand.Err("unknown command")
		return nil, fmt.Errorf("unknown command: %s", rawCommand.Command)
	}

	cmd, err := cmdType.FromRaw(rawCommand)
	if err != nil {
		return nil, fmt.Errorf("from raw: %w", err)
	}

	return cmd, nil
}

type CommandDecoder struct {
	wsconn  *Conn
	palette CommandPalette
}
