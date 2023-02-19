package chat_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/irth/ovencast/chat/internal/chat"
	"github.com/irth/wsrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	chat   *chat.Chat
	server *httptest.Server
}

func NewTestServer(t *testing.T) *TestServer {
	chat, err := chat.NewChat()
	require.NoError(t, err)
	s := httptest.NewServer(chat.Handler())

	return &TestServer{
		chat:   chat,
		server: s,
	}
}

func (ts *TestServer) Close() {
	ts.server.Close()
}

func (ts *TestServer) Connect(t *testing.T) TestConn {
	u := "ws" + strings.TrimPrefix(ts.server.URL, "http") + "/ws"
	conn, _, _, err := ws.Dial(context.Background(), u)
	require.NoError(t, err)
	return NewTestConn(conn)
}

type TestConn struct {
	c net.Conn

	r *wsutil.Reader
	w *wsutil.Writer

	d *json.Decoder
	e *json.Encoder
}

func NewTestConn(conn net.Conn) TestConn {
	r := wsutil.NewClientSideReader(conn)
	w := wsutil.NewWriter(conn, ws.StateClientSide, ws.OpText)
	return TestConn{
		c: conn,

		r: r,
		w: w,

		d: json.NewDecoder(r),
		e: json.NewEncoder(w),
	}
}

func (tc TestConn) Close() {
	tc.c.Close()
}

func (tc TestConn) Read(t *testing.T, obj any) bool {
	hdr, err := tc.r.NextFrame()
	require.NoError(t, err)

	if hdr.OpCode == ws.OpClose {
		return false
	}

	err = tc.d.Decode(obj)
	require.NoError(t, err)

	return true
}

func (tc TestConn) Write(t *testing.T, obj any) {
	err := tc.e.Encode(obj)
	require.NoError(t, err)

	err = tc.w.Flush()
	require.NoError(t, err)
}

func (tc TestConn) Call(t *testing.T, req any, res any) {
	tc.Write(t, req)
	tc.Read(t, res)
}

func TestSetNick(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	conn := ts.Connect(t)
	defer conn.Close()

	var response wsrpc.Response[any]
	conn.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname1",
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
		"request": "testNickname1",
	}, &response)
	assert.True(t, response.OK)

	var response2 wsrpc.Response[any]
	conn.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname1",
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
		"request": "testNickname1",
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to set a nickname that is not in use")

	conn2 := ts.Connect(t)
	defer conn2.Close()

	var response2 wsrpc.Response[any]
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname1",
	}, &response2)
	assert.False(t, response2.OK, "it should not be possible to use a nickname that another connection is already using")

	response2.OK = false
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname2",
	}, &response2)
	assert.True(t, response2.OK, "it should be possible to set a nickname that's not in use even after a previous failure to do so")

	response1.OK = false
	conn1.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname3",
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to change a nickname to one that is not in use")

	response2.OK = false
	conn2.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname1",
	}, &response2)
	assert.True(t, response2.OK, "it should be possible to reuse a nickname after it's been let go")

	response1.OK = false
	conn2.Close()
	<-time.After(100 * time.Millisecond) // wait for the server to process the disconnection
	conn1.Call(t, map[string]any{
		"command": "nick",
		"request": "testNickname1",
	}, &response1)
	assert.True(t, response1.OK, "it should be possible to use a nickname after it's original user disconnects")
}
