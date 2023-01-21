package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WSConn struct {
	w *wsutil.Writer
	r *wsutil.Reader

	encoder *json.Encoder
	decoder *json.Decoder
}

func NewWSConn(conn net.Conn) WSConn {
	w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	r := wsutil.NewReader(conn, ws.StateServerSide)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r)

	return WSConn{
		w: w, r: r,
		encoder: encoder,
		decoder: decoder,
	}
}

type WSCommand[T any, U any] struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Request T      `json:"request"`

	wsconn WSConn
}

type WSCommandInterface interface {
	Error(string, ...interface{}) error
}

func (w WSCommand[T, U]) json(obj interface{}) error {
	err := w.wsconn.encoder.Encode(obj)
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	if err = w.wsconn.w.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}

func (w WSCommand[T, U]) Reply(reply U) error {
	response := WSResponse[U]{
		ID:       w.ID,
		OK:       true,
		Command:  w.Command,
		Response: reply,
	}

	return w.json(response)
}

func (w WSCommand[T, U]) Error(errMsg string, args ...interface{}) error {
	response := WSResponse[Empty]{
		ID:      w.ID,
		OK:      false,
		Command: w.Command,
		Error:   fmt.Sprintf(errMsg, args...),
	}

	return w.json(response)
}

type WSResponse[T any] struct {
	ID       string `json:"id"`
	OK       bool   `json:"ok"`
	Command  string `json:"command"`
	Response T      `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (w WSConn) DecodeCommand() (ret WSCommandInterface, err error) {
	rawCommand := WSCommand[json.RawMessage, Empty]{
		wsconn: w,
	}

	hdr, err := w.r.NextFrame()
	if err != nil {
		return nil, fmt.Errorf("ws NextFrame: %w", err)
	}
	if hdr.OpCode == ws.OpClose {
		return nil, io.EOF
	}

	err = w.decoder.Decode(&rawCommand)
	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	fmt.Printf("%+v", rawCommand)

	switch rawCommand.Command {
	case "ping":
		cmd := Ping{
			ID:     rawCommand.ID,
			wsconn: w,
		}
		err = json.Unmarshal(rawCommand.Request, &cmd.Request)
		ret = cmd
	default:
		err = fmt.Errorf("unknown command")
		rawCommand.Error(err.Error())
	}

	return
}
