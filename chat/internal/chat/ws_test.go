package chat_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/irth/ovencast/chat/internal/chat"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	chat   *chat.Chat
	server *httptest.Server
}

func NewTestServer(t *testing.T) *TestServer {
	chat, err := chat.NewChat()
	require.NoError(t, err)
	go chat.Start(context.Background())
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

	if chat.DEBUG {
		fmt.Printf("%s<<< ", Cyan)
		json.NewEncoder(os.Stdout).Encode(obj)
		fmt.Printf("%s", Reset)
	}
	return true
}

var Blue = "\033[34m"
var Cyan = "\033[36m"
var Reset = "\033[0m"

func (tc TestConn) Write(t *testing.T, obj any) {
	err := tc.e.Encode(obj)
	require.NoError(t, err)

	err = tc.w.Flush()
	require.NoError(t, err)

	if chat.DEBUG {
		fmt.Printf("%s>>> ", Blue)
		json.NewEncoder(os.Stdout).Encode(obj)
		fmt.Printf("%s", Reset)
	}
}

func (tc TestConn) Call(t *testing.T, req any, res any) {
	tc.Write(t, req)
	tc.Read(t, res)
	fmt.Println("")
}
