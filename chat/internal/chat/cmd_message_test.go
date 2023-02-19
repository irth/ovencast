package chat_test

import (
	"testing"

	"github.com/irth/wsrpc"
	"github.com/stretchr/testify/require"
)

type j = map[string]any

func cmd(name string, req any) j {
	return j{
		"command": name,
		"request": req,
	}
}

func TestCannotSendMessageWithoutANick(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn1 := ts.Connect(t)
	defer conn1.Close()

	var response wsrpc.Response[any]
	conn1.Call(t, cmd("message", j{"content": "this is the message content"}), &response)
	require.False(t, response.OK, "sending a message with no nick set should fail")

	conn1.Call(t, cmd("nick", j{"nickname": "ANickname"}), &response)
	require.True(t, response.OK)

	response.OK = false
	conn1.Call(t, cmd("message", j{"content": "this is the message content"}), &response)
	require.True(t, response.OK, "sending a message with a nick set should succeed")
}

func TestMessagesGetBroadcasted(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn1 := ts.Connect(t)
	defer conn1.Close()

	conn2 := ts.Connect(t)
	defer conn2.Close()

	conn3 := ts.Connect(t)
	defer conn3.Close()

	nick := "ANickname"
	content := "this is a test message"

	var response wsrpc.Response[any]
	conn1.Call(t, cmd("nick", j{"nickname": nick}), &response)
	require.True(t, response.OK)

	conn1.Call(t, cmd("message", j{"content": content}), &response)
	require.True(t, response.OK, "sending a message with a nick set should succeed")

	var m1 struct {
		Type    string `json:"type"`
		Message struct {
			Nickname string `json:"nickname"`
			Content  string `json:"content"`
		} `json:"message"`
	}
	m2 := m1

	conn2.Read(t, &m1)
	require.Equal(t, "message", m1.Type)
	require.Equal(t, nick, m1.Message.Nickname)
	require.Equal(t, content, m1.Message.Content)

	conn3.Read(t, &m2)
	require.Equal(t, "message", m2.Type)
	require.Equal(t, nick, m2.Message.Nickname)
	require.Equal(t, content, m2.Message.Content)
}

func TestNickChangesAreReflectedInMessages(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn1 := ts.Connect(t)
	defer conn1.Close()

	nick1 := "testNickname1"
	nick2 := "testNickname2"

	var r wsrpc.Response[any]
	conn1.Call(t, cmd("nick", j{"nickname": nick1}), &r)
	require.True(t, r.OK)

	conn1.Call(t, cmd("message", j{"content": "message content"}), &r)
	require.True(t, r.OK)

	var msg j
	conn1.Read(t, &msg)
	require.Equal(t, "message", msg["type"])
	require.Equal(t, nick1, msg["message"].(j)["nickname"])

	conn1.Call(t, cmd("nick", j{"nickname": nick2}), &r)
	require.True(t, r.OK)

	conn1.Call(t, cmd("message", j{"content": "message content"}), &r)
	require.True(t, r.OK)

	conn1.Read(t, &msg)
	require.Equal(t, "message", msg["type"])
	require.Equal(t, nick2, msg["message"].(j)["nickname"])
}
