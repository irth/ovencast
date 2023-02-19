package chat_test

import (
	"testing"
	"time"

	"github.com/irth/wsrpc"
	"github.com/stretchr/testify/assert"
)

func TestSetNick(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn := ts.Connect(t)
	defer conn.Close()

	var response wsrpc.Response[any]
	conn.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response)

	assert.True(t, response.OK)
}

func TestSetSameNickTwiceSucceeds(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn := ts.Connect(t)
	defer conn.Close()

	var response wsrpc.Response[any]
	conn.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response)
	assert.True(t, response.OK)

	var response2 wsrpc.Response[any]
	conn.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response2)
}

func TestNicknameReservation(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn1 := ts.Connect(t)
	defer conn1.Close()

	var response1 wsrpc.Response[any]
	conn1.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to set a nickname that is not in use")

	conn2 := ts.Connect(t)
	defer conn2.Close()

	var response2 wsrpc.Response[any]
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response2)
	assert.False(t, response2.OK, "it should not be possible to use a nickname that another connection is already using")

	response2.OK = false
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname2"},
	}, &response2)
	assert.True(t, response2.OK, "it should be possible to set a nickname that's not in use even after a previous failure to do so")

	response1.OK = false
	conn1.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname3"},
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to change a nickname to one that is not in use")

	response2.OK = false
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response2)
	assert.True(t, response2.OK, "it should be possible to reuse a nickname after it's been let go")

	response1.OK = false
	conn2.Close()
	<-time.After(100 * time.Millisecond) // wait for the server to process the disconnection
	conn1.Call(t, map[string]any{
		"command": "nick",
		"request": map[string]any{"nickname": "testNickname1"},
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to use a nickname after it's original user disconnects")
}
